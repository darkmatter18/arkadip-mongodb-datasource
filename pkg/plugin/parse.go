package plugin

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ExtractPartsOfMongoCommand(command string) ([]string, error) {
	// Remove /t, /n and spaces and make the string single line
	singlelineStr := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(command, "\t", ""), "\n", ""), " ", "")

	// Split the command string into parts
	regex := regexp.MustCompile(`(?s)([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)\.(\w+)\((.*?)\)\s*;?\s*$`)
	matches := regex.FindStringSubmatch(singlelineStr)
	if matches == nil {
		return nil, fmt.Errorf("regex parsing Error : %s", command)
	}

	return matches[1:], nil
}

func MongoPipeline(str string) (mongo.Pipeline, error) {
	var pipeline = []bson.D{}
	var err error
	str = strings.TrimSpace(str)
	if strings.Index(str, "[") != 0 {
		var doc bson.D
		err = bson.UnmarshalExtJSON([]byte(str), false, &doc)
		pipeline = append(pipeline, doc)
	} else {
		err = bson.UnmarshalExtJSON([]byte(str), false, &pipeline)
	}
	return pipeline, err
}

func MongoFind(str string) (bson.D, error) {
	var pipeline = bson.D{}
	str = strings.TrimSpace(str)
	err := bson.UnmarshalExtJSON([]byte(str), false, &pipeline)
	return pipeline, err
}

func MongoQuery(client *mongo.Client, ctx context.Context, db string, collection string, operation string, query string) (*mongo.Cursor, error) {
	db_collection := client.Database(db).Collection(collection)

	switch operation {
	case "find":
		pipeline, err := MongoFind(query)
		if err != nil {
			return nil, err
		}
		return db_collection.Find(ctx, pipeline)
	case "aggregate":
		pipeline, err := MongoPipeline(query)
		if err != nil {
			return nil, err
		}
		return db_collection.Aggregate(ctx, pipeline)
	default:
		return nil, fmt.Errorf("following option %s in not available", operation)
	}
}
