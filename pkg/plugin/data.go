package plugin

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func purseAnyToString(value interface{}) string {
	var data string = ""
	switch i := value.(type) {
	case int:
		data = strconv.FormatInt(int64(i), 10)
	case int32:
		data = strconv.FormatInt(int64(i), 10)
	case int64:
		data = strconv.FormatInt(i, 10)
	case float32:
		data = strconv.FormatFloat(float64(i), 'f', -1, 32)
	case float64:
		data = strconv.FormatFloat(i, 'f', -1, 32)
	case bool:
		data = strconv.FormatBool(i)
	case string:
		data = i
	case primitive.ObjectID:
		data = i.Hex()
	case primitive.Binary: // UUID
		if u, e := uuid.FromBytes(i.Data); e != nil {
			fmt.Print(e)
		} else {
			data = u.String()
		}
	case primitive.DateTime:
		data = i.Time().Format(time.RFC3339)
	case primitive.M:
		data = purseMapToString(i)
	case primitive.A:
		data = purseArrayToString(i)
	case nil:
		fmt.Printf(" x is nil")
	default:
		fmt.Printf(" don't know the type ")
		fmt.Println("%V %V", reflect.TypeOf(value), i)
	}

	return data
}

func purseMapToString(value bson.M) string {
	rval := "{"
	for key, value := range value {
		content := purseAnyToString(value)
		rval += fmt.Sprintf(`"%s": "%s",`, key, content)
	}
	rval += "}"
	return rval
}

func purseArrayToString(value bson.A) string {
	rval := "["
	for _, value := range value {
		content := purseAnyToString(value)
		rval += fmt.Sprintf(`"%s",`, content)
	}
	rval += "]"
	return rval
}
