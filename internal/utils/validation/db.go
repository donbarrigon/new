package validation

import (
	"context"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/fm"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// IsUnique validates that a value is unique in the specified collection and field
// Params: [collection, field, excludeId (optional)]
func Unique(value reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	if len(params) < 2 {
		return "Parámetro colección y campo son requeridos", fm.Placeholder{}, true
	}

	col := params[0]
	field := params[1]
	ph := fm.Placeholder{"col": col, "field": field}

	// Skip validation if value is nil or empty
	if !value.IsValid() || value.IsZero() {
		return "", ph, true
	}

	// Build the filter
	filter := bson.M{field: value.Interface()}

	// If excludeId is provided (for updates), exclude that document from the check
	if len(params) >= 3 && params[2] != "" {
		excludeId := params[2]

		// Try to convert to ObjectID if it looks like one
		if oid, err := primitive.ObjectIDFromHex(excludeId); err == nil {
			filter["_id"] = bson.M{"$ne": oid}
		} else {
			filter["_id"] = bson.M{"$ne": excludeId}
		}

		ph["excludeId"] = excludeId
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
			return "", ph, true
		}

		// Other errors (connection, timeout, etc.)
		return "Error al validar unicidad: " + err.Error(), ph, true
	}

	// Document found, value is not unique
	ph["value"] = value.Interface()
	return "El valor ya existe en la colección", ph, false
}

// Exists validates that a value exists in the specified collection and field
// Params: [collection, field]
func Exists(value reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	if len(params) < 2 {
		return "Parámetro colección y campo son requeridos", fm.Placeholder{}, true
	}

	col := params[0]
	field := params[1]
	ph := fm.Placeholder{"col": col, "field": field}

	// Skip validation if value is nil or empty
	if !value.IsValid() || value.IsZero() {
		return "El valor es requerido", ph, false
	}

	// Build the filter
	filter := bson.M{field: value.Interface()}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if document exists
	var result map[string]any
	err := db.Col(col).FindOne(ctx, filter).Decode(&result)

	if err != nil {
		// If no document found, value doesn't exist
		if err == mongo.ErrNoDocuments {
			ph["value"] = value.Interface()
			return "El valor no existe en la colección", ph, false
		}

		// Other errors (connection, timeout, etc.)
		return "Error al validar existencia: " + err.Error(), ph, false
	}

	// Document found, value exists
	return "", ph, true
}

// ExistsById validates that a document exists by its ID
// Params: [collection]
func ExistsById(value reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	if len(params) < 1 {
		return "Parámetro colección es requerido", fm.Placeholder{}, true
	}

	col := params[0]
	ph := fm.Placeholder{"col": col}

	// Skip validation if value is nil or empty
	if !value.IsValid() || value.IsZero() {
		return "El ID es requerido", ph, false
	}

	var filter bson.M

	// Handle different ID types
	switch v := value.Interface().(type) {
	case string:
		// Try to convert string to ObjectID
		if oid, err := primitive.ObjectIDFromHex(v); err == nil {
			filter = bson.M{"_id": oid}
		} else {
			filter = bson.M{"_id": v}
		}
	case primitive.ObjectID:
		filter = bson.M{"_id": v}
	default:
		filter = bson.M{"_id": v}
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if document exists
	var result map[string]any
	err := db.Col(col).FindOne(ctx, filter).Decode(&result)

	if err != nil {
		// If no document found, ID doesn't exist
		if err == mongo.ErrNoDocuments {
			ph["id"] = value.Interface()
			return "El documento no existe", ph, false
		}

		// Other errors (connection, timeout, etc.)
		return "Error al validar existencia del documento: " + err.Error(), ph, false
	}

	// Document found, ID exists
	return "", ph, true
}

// UniqueComposite validates that a combination of fields is unique
// Params: [collection, field1, field2, ..., excludeId (optional if last param starts with "exclude:")]
func UniqueComposite(values map[string]reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	if len(params) < 2 {
		return "Se requiere al menos colección y un campo", fm.Placeholder{}, true
	}

	col := params[0]
	fields := params[1:]
	ph := fm.Placeholder{"col": col}

	var excludeId string

	// Check if last param is excludeId
	if len(fields) > 0 && len(fields[len(fields)-1]) > 8 && fields[len(fields)-1][:8] == "exclude:" {
		excludeId = fields[len(fields)-1][8:] // Remove "exclude:" prefix
		fields = fields[:len(fields)-1]       // Remove excludeId from fields
		ph["excludeId"] = excludeId
	}

	// Build the filter with all field values
	filter := bson.M{}
	for _, field := range fields {
		if value, exists := values[field]; exists && value.IsValid() && !value.IsZero() {
			filter[field] = value.Interface()
			ph[field] = value.Interface()
		}
	}

	// If no valid fields to check, pass validation
	if len(filter) == 0 {
		return "", ph, true
	}

	// Add exclude condition if provided
	if excludeId != "" {
		if oid, err := primitive.ObjectIDFromHex(excludeId); err == nil {
			filter["_id"] = bson.M{"$ne": oid}
		} else {
			filter["_id"] = bson.M{"$ne": excludeId}
		}
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if combination exists
	var result map[string]any
	err := db.Col(col).FindOne(ctx, filter).Decode(&result)

	if err != nil {
		// If no document found, combination is unique
		if err == mongo.ErrNoDocuments {
			return "", ph, true
		}

		// Other errors
		return "Error al validar unicidad compuesta: " + err.Error(), ph, true
	}

	// Document found, combination is not unique
	return "La combinación de valores ya existe", ph, false
}
