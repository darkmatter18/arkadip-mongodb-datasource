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

func MongoPipeline(str string) mongo.Pipeline {
	var pipeline = []bson.D{}
	str = strings.TrimSpace(str)
	if strings.Index(str, "[") != 0 {
		var doc bson.D
		bson.UnmarshalExtJSON([]byte(str), false, &doc)
		pipeline = append(pipeline, doc)
	} else {
		bson.UnmarshalExtJSON([]byte(str), false, &pipeline)
	}
	return pipeline
}

func MongoFind(str string) bson.D {
	var pipeline = bson.D{}
	str = strings.TrimSpace(str)
	bson.UnmarshalExtJSON([]byte(str), false, &pipeline)
	return pipeline
}

func MongoQuery(client *mongo.Client, ctx context.Context, db string, collection string, operation string, query string) (*mongo.Cursor, error) {
	db_collection := client.Database(db).Collection(collection)

	switch operation {
	case "find":
		return db_collection.Find(ctx, MongoFind(query))
	case "aggregate":
		return db_collection.Aggregate(ctx, MongoPipeline(query))
	default:
		return nil, fmt.Errorf("following option %s in not available", operation)
	}
}
