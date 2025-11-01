package model

import (
	"context"
	"donbarrigon/new/internal/app/data/validator"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/err"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role struct {
	ID          bson.ObjectID `bson:"_id"         json:"id"`
	Name        string        `bson:"name"        json:"name"`
	Permissions []string      `bson:"permissions" json:"permissions"`
}

// ================================================================
//                  FUNCIONES PARA LA INTERFAZ
// ================================================================

func (r *Role) GetID() bson.ObjectID { return r.ID }
func (r *Role) Coll() string         { return "roles" }

// ================================================================
//                  FUNCIONES CRUD
// ================================================================

func RoleByID(id bson.ObjectID) (*Role, error) {
	role := &Role{}
	filter := bson.D{{Key: "_id", Value: id}}
	if e := db.Mongo.Collection(role.Coll()).FindOne(context.TODO(), filter).Decode(role); e != nil {
		return nil, err.Mongo(e)
	}
	return role, nil
}

func RoleByHexID(id string) (*Role, error) {
	oid, e := bson.ObjectIDFromHex(id)
	if e != nil {
		return nil, err.HexID(e.Error())
	}
	return RoleByID(oid)
}

func RoleByName(name string) (*Role, error) {
	role := &Role{}
	filter := bson.D{{Key: "name", Value: name}}
	if e := db.Mongo.Collection(role.Coll()).FindOne(context.TODO(), filter).Decode(role); e != nil {
		return nil, err.Mongo(e)
	}
	return role, nil
}

func RoleCreate(dto *validator.RoleStore) (*Role, error) {
	role := &Role{
		Name:        dto.Name,
		Permissions: []string{},
	}
	result, e := db.Mongo.Collection(role.Coll()).InsertOne(context.TODO(), role)
	if e != nil {
		return nil, err.Mongo(e)
	}
	role.ID = result.InsertedID.(bson.ObjectID)
	return role, nil
}

func (role *Role) Update(dto *validator.RoleStore) (*Changes, error) {

	changes := NewChanges()
	if role.Name != dto.Name {
		changes.Old["name"] = role.Name
		changes.New["name"] = dto.Name
		role.Name = dto.Name
	}

	filter := bson.D{{Key: "_id", Value: role.ID}}
	result, e := db.Mongo.Collection(role.Coll()).UpdateOne(context.TODO(), filter, role)
	if e != nil {
		return nil, err.Mongo(e)
	}
	return changes, err.MongoUpdateResult(result)
}
