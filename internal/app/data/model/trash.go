package model

import (
	"context"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/err"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Trash struct {
	ID         bson.ObjectID `bson:"_id"           json:"id"`
	UserID     bson.ObjectID `bson:"user_id"       json:"userId"`
	Collection string        `bson:"collection"    json:"collection"`
	Data       any           `bson:"data"          json:"data"`
	DeletedAt  time.Time     `bson:"deleted_at"    json:"deletedAt"`
}

// ================================================================
//                 FUNCIONES PARA LA INTERFAZ
// ================================================================

func (t *Trash) GetID() bson.ObjectID { return t.ID }
func (t *Trash) Coll() string         { return "trash" }

// ================================================================
//                    FUNCIONES AUXILIARES
// ================================================================

func MoveToTrash(userID bson.ObjectID, m db.MongoModel) error {
	trash := &Trash{
		UserID:     userID,
		Collection: m.Coll(),
		Data:       m,
		DeletedAt:  time.Now(),
	}

	_, e := db.Mongo.Collection("trash").InsertOne(context.TODO(), trash)
	if e != nil {
		return err.Mongo(e)
	}

	filter := bson.D{bson.E{Key: "_id", Value: m.GetID()}}

	result, e := db.Mongo.Collection(m.Coll()).DeleteOne(context.TODO(), filter)
	if e != nil {
		return err.Mongo(e)
	}
	return err.MongoDeleteResult(result)
}

func RestoreByID(m db.MongoModel, coll string, oid bson.ObjectID) error {
	filter := bson.D{{Key: "data._id", Value: oid}, {Key: "collection", Value: coll}}
	if e := db.Mongo.Collection("trash").FindOne(context.TODO(), filter).Decode(m); e != nil {
		return err.Mongo(e)
	}

	if _, e := db.Mongo.Collection(coll).InsertOne(context.TODO(), m); e != nil {
		return err.Mongo(e)
	}

	result, e := db.Mongo.Collection("trash").DeleteOne(context.TODO(), filter)
	if e != nil {
		return err.Mongo(e)
	}

	return err.MongoDeleteResult(result)
}

func RestoreByHexID(m db.MongoModel, coll string, id string) error {
	oid, e := bson.ObjectIDFromHex(id)
	if e != nil {
		return err.HexID(e)
	}
	return RestoreByID(m, coll, oid)
}
