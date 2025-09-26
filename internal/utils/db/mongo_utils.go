package db

import (
	"context"
	"donbarrigon/new/internal/utils/config"
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/logs"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitMongoDB() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(config.DbConnectionString).SetServerAPIOptions(serverAPI)
	clientOptions.SetMaxPoolSize(100)
	clientOptions.SetMinPoolSize(5)
	clientOptions.SetRetryWrites(true)
	clientOptions.SetTimeout(30 * time.Second)

	var err error

	DBClient, err = mongo.Connect(clientOptions)
	if err != nil {
		logs.Info("🔴💥 Fail to connect db %s: %s", config.DbName, config.DbConnectionString)
		return err
	}
	DB = DBClient.Database(config.DbName)

	logs.Info("🍃 Successful connection to %s: %s", config.DbName, config.DbConnectionString)
	return nil
}

func CloseMongoDB() error {

	if DBClient == nil {
		return nil
	}

	return DBClient.Disconnect(context.TODO())
}

// Fill compara model con request y retorna los valores antiguos y los nuevos.
// Usa el tag bson como clave.
// Además actualiza el model con los valores nuevos.
// @return original, dirty, error
func Fill(model any, request any) (map[string]any, map[string]any, err.Error) {

	original := map[string]any{}
	dirty := map[string]any{}

	modelValue := reflect.ValueOf(model)
	requestValue := reflect.ValueOf(request)

	if modelValue.Kind() != reflect.Ptr || requestValue.Kind() != reflect.Ptr {
		return original, dirty, err.Internal("The parameters model and request must be pointers")
	}

	modelValue = modelValue.Elem()
	requestValue = requestValue.Elem()

	if modelValue.Kind() != reflect.Struct || requestValue.Kind() != reflect.Struct {
		return original, dirty, err.Internal("The parameters model and request must be structs")
	}

	modelType := modelValue.Type()
	requestType := requestValue.Type()

	for i := 0; i < requestType.NumField(); i++ {
		requestField := requestType.Field(i)
		fieldName := requestField.Name
		requestFieldValue := requestValue.Field(i)

		if !requestField.IsExported() {
			continue
		}

		modelField := modelValue.FieldByName(fieldName)
		if !modelField.IsValid() || !modelField.CanSet() {
			continue
		}

		if requestFieldValue.IsZero() {
			continue
		}

		var newValue reflect.Value
		if modelField.Type() == reflect.TypeOf(bson.ObjectID{}) &&
			requestFieldValue.Kind() == reflect.String {

			oid, err := bson.ObjectIDFromHex(requestFieldValue.String())
			if err != nil {
				continue
			}
			newValue = reflect.ValueOf(oid)

		} else if requestFieldValue.Type().AssignableTo(modelField.Type()) {
			newValue = requestFieldValue
		} else if requestFieldValue.Type().ConvertibleTo(modelField.Type()) {
			newValue = requestFieldValue.Convert(modelField.Type())
		} else {
			continue
		}

		currentValue := modelField
		if !reflect.DeepEqual(currentValue.Interface(), newValue.Interface()) {
			// sacar el tag bson
			modelFieldStruct, found := modelType.FieldByName(fieldName)
			if !found {
				continue
			}

			bsonTag := modelFieldStruct.Tag.Get("bson")
			if bsonTag == "" {
				bsonTag = strings.ToLower(fieldName)
			} else if commaIndex := strings.Index(bsonTag, ","); commaIndex != -1 {
				bsonTag = bsonTag[:commaIndex]
			}

			// Guardamos viejo y nuevo
			original[bsonTag] = currentValue.Interface()
			dirty[bsonTag] = newValue.Interface()

			// Actualizamos el modelo
			modelField.Set(newValue)
		}
	}

	return original, dirty, nil
}
