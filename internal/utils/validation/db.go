package validation

import (
	"context"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/fm"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Unique validates that a value is unique in the specified collection and field
// Params: [collection, field, excludeId (optional)]
func IsUnique(value any, col string, field string, excludeId string) (string, fm.Placeholder, bool) {

	ph := fm.Placeholder{}

	// Build the filter
	filter := bson.M{field: value}

	// Try to convert to ObjectID if it looks like one
	if oid, err := bson.ObjectIDFromHex(excludeId); err == nil {
		filter["_id"] = bson.M{"$ne": oid}
	} else {
		filter["_id"] = bson.M{"$ne": excludeId}
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if document exists
	var result map[string]any
	err := db.Col(col).FindOne(ctx, filter).Decode(&result)

	if err != nil {
		// If no document found, value is unique
		if err == mongo.ErrNoDocuments {
			return "", nil, false
		}

		// Other errors (connection, timeout, etc.)
		return "Error al validar unicidad del campo :field:" + err.Error(), ph, true
	}

	// Document found, value is not unique
	// ph["value"] = value.Interface()
	return "El :field ya existe", ph, true
}

// Exists validates that a value exists in the specified collection and field
// Params: [collection, field]
func Exists(value reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	if len(params) < 2 {
		return "Parámetro colección y campo son requeridos", fm.Placeholder{}, true
	}

	col := params[0]
	field := params[1]
	ph := fm.Placeholder{}

	// Skip validation if value is nil or empty
	if !value.IsValid() || value.IsZero() {
		return "El :field es requerido", ph, true
	}

	// Build the filter
	var filter bson.M
	if field == "_id" {
		switch v := value.Interface().(type) {
		case string:
			// Try to convert string to ObjectID
			if oid, err := bson.ObjectIDFromHex(v); err == nil {
				filter = bson.M{"_id": oid}
			} else {
				filter = bson.M{"_id": v}
			}
		case bson.ObjectID:
			filter = bson.M{"_id": v}
		default:
			filter = bson.M{"_id": v}
		}
	} else {
		filter = bson.M{field: value.Interface()}
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if document exists
	var result map[string]any
	err := db.Col(col).FindOne(ctx, filter).Decode(&result)

	if err != nil {
		// If no document found, value doesn't exist
		if err == mongo.ErrNoDocuments {
			// ph["value"] = fmt.Sprintf("%v", value.Interface())
			return "El :field no existe", ph, true
		}

		// Other errors (connection, timeout, etc.)
		return "Error al validar existencia del campo :field: " + err.Error(), ph, true
	}

	// Document found, value exists
	return "", nil, false
}
