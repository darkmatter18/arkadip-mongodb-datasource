package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces- only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

// NewDatasource creates a new datasource instance.
func NewDatasource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {

	// Get db_uri from JSON data
	db_uri, exists := settings.DecryptedSecureJSONData["db_uri"]
	if !exists {
		// Use the decrypted API key.
		return nil, fmt.Errorf("not able to found db uri")
	}
	log.Print(db_uri)

	// Get JSON data
	var data map[string]interface{}
	ere := json.Unmarshal(settings.JSONData, &data)

	if ere != nil {
		return nil, ere
	}
	test_db := data["test_db"].(string)
	log.Print(test_db)

	// Mongo DB Connection
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(db_uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		return nil, err
	}

	// Testing of the DB
	var result bson.M
	if err := client.Database(test_db).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return nil, err
	}

	return &Datasource{
		db:      *client,
		test_db: test_db,
	}, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct {
	db      mongo.Client
	test_db string
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
	err := d.db.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	Q string `json:"q"`
}

func (d *Datasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	// Unmarshal the JSON into our queryModel.
	var qm queryModel

	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}

	log.Printf(qm.Q)

	if len(qm.Q) == 0 {
		var response backend.DataResponse
		return response
	}

	cmd, parts_err := ExtractPartsOfMongoCommand(qm.Q)
	if parts_err != nil {
		var response backend.DataResponse
		response.Error = parts_err
		return response
	}

	mongo_data, query_err := MongoQuery(&d.db, context.TODO(), cmd[0], cmd[1], cmd[2], cmd[3])
	if query_err != nil {
		var response backend.DataResponse
		response.Error = query_err
		return response
	}

	// defer mongo_data.Close(context.Background())

	out_data := make(map[string]interface{})
	length := 0

	for mongo_data.Next(context.TODO()) {
		var m bson.M
		if err := mongo_data.Decode(&m); err != nil {
			fmt.Printf("Error in JSON")
		}

		for key, value := range m {
			var data string = ""
			// fmt.Println("%V %V", key, reflect.TypeOf(value))
			switch i := value.(type) {
			case primitive.DateTime:
				if v, ok := out_data[key].([]time.Time); ok {
					out_data[key] = append(v, i.Time())
				} else {
					out_data[key] = []time.Time{i.Time()}
				}
				continue
			default:
				data = purseAnyToString(i)
			}

			if v, ok := out_data[key].([]string); ok {
				if length > len(v) {
					padding_length := length - len(v)
					for i := 0; i < padding_length; i++ {
						v = append(v, "")
					}
				}
				out_data[key] = append(v, data)
			} else {
				// Initial Padding
				// If there is no value in the beginning but other keys have values
				// then fill the current key value with a array of len and fill with 0
				if length > 0 {
					d := make([]string, length)
					out_data[key] = append(d, data)
				} else {
					out_data[key] = []string{data}
				}
			}
		}
		length++
	}

	// create data frame response.
	// For an overview on data frames and how grafana handles them:
	// https://grafana.com/docs/grafana/latest/developers/plugins/data-frames/
	frame := data.NewFrame("response")

	for key, value := range out_data {
		// add fields.
		frame.Fields = append(frame.Fields,
			data.NewField(key, nil, value),
		)
	}

	// add the frames to the response.
	var response backend.DataResponse
	response.Frames = append(response.Frames, frame)

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Data source is working"

	// Testing of the DB
	var result bson.M
	if err := d.db.Database(d.test_db).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		status = backend.HealthStatusError
		message = err.Error()

		return &backend.CheckHealthResult{
			Status:  status,
			Message: message,
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
