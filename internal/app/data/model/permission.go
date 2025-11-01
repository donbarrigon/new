package model

import (
	"context"
	"donbarrigon/new/internal/app/data/validator"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/err"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Permission struct {
	ID   bson.ObjectID `bson:"_id"  json:"id"`
	Name string        `bson:"name" json:"name"`
}

// ================================================================
//                  FUNCIONES PARA LA INTERFAZ
// ================================================================

func (p *Permission) GetID() bson.ObjectID { return p.ID }
func (p *Permission) Coll() string         { return "permissions" }

// ================================================================
//                      FUNCIONES CRUD
// ================================================================

func PermissionByID(id bson.ObjectID) (*Permission, error) {
	permission := &Permission{}
	filter := bson.D{{Key: "_id", Value: id}}
	if e := db.Mongo.Collection(permission.Coll()).FindOne(context.TODO(), filter).Decode(permission); e != nil {
		return nil, err.Mongo(e)
	}
	return permission, nil
}

func PermissionByHexID(id string) (*Permission, error) {
	oid, e := bson.ObjectIDFromHex(id)
	if e != nil {
		return nil, err.HexID(e.Error())
	}
	return PermissionByID(oid)
}

func PermissionByName(name string) (*Permission, error) {
	permission := &Permission{}
	filter := bson.D{{Key: "name", Value: name}}
	if e := db.Mongo.Collection(permission.Coll()).FindOne(context.TODO(), filter).Decode(permission); e != nil {
		return nil, err.Mongo(e)
	}
	return permission, nil
}

func PermissionCreate(name string) (*Permission, error) {
	permission := &Permission{
		Name: name,
	}
	result, e := db.Mongo.Collection(permission.Coll()).InsertOne(context.TODO(), permission)
	if e != nil {
		return nil, err.Mongo(e)
	}
	permission.ID = result.InsertedID.(bson.ObjectID)
	return permission, nil
}

func (permission *Permission) Update(dto *validator.PermissionStore) (*Changes, error) {

	changes := NewChanges()
	if permission.Name != dto.Name {
		changes.Old["name"] = permission.Name
		changes.New["name"] = dto.Name
		permission.Name = dto.Name
	}

	filter := bson.D{{Key: "_id", Value: permission.ID}}
	update := bson.D{{Key: "$set", Value: permission}}
	result, e := db.Mongo.Collection(permission.Coll()).UpdateOne(context.TODO(), filter, update)
	if e != nil {
		return nil, err.Mongo(e)
	}
	if e := err.MongoUpdateResult(result); e != nil {
		return nil, e
	}

	// actualizo los roles
	filter = bson.D{}
	update = bson.D{{Key: "$set", Value: bson.D{{Key: "permissions.$[elem]", Value: permission.Name}}}}
	opts := options.UpdateMany().SetArrayFilters([]any{
		bson.D{{Key: "elem", Value: changes.Old["name"]}},
	})
	resultr, e := db.Mongo.Collection("roles").UpdateMany(context.TODO(), filter, update, opts)
	if e != nil {
		return nil, err.Mongo(e)
	}
	if e := err.MongoUpdateResult(resultr); e != nil {
		return nil, e
	}

	// actualizo los usuarios

	return changes, nil
}
