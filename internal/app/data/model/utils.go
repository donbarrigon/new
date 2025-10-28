package model

import (
	"context"
	"donbarrigon/new/internal/utils/auth"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/err"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ================================================================
// trash
// ================================================================
type Trash struct {
	ID         bson.ObjectID `bson:"_id"           json:"id"`
	UserID     bson.ObjectID `bson:"user_id"       json:"userId"`
	Collection string        `bson:"collection"    json:"collection"`
	Data       any           `bson:"data"          json:"data"`
	DeletedAt  time.Time     `bson:"deleted_at"    json:"deletedAt"`
}

func (t *Trash) GetID() bson.ObjectID { return t.ID }
func (t *Trash) Coll() string         { return "trash" }

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

// ================================================================
// History
// ================================================================

type History struct {
	ID           bson.ObjectID  `bson:"_id"           json:"id"`
	UserID       bson.ObjectID  `bson:"user_id"       json:"userId"`
	CollectionID bson.ObjectID  `bson:"collection_id" json:"collectionId"`
	Collection   string         `bson:"collection"    json:"collection"`
	OldData      map[string]any `bson:"old_data"      json:"oldData"`
	Action       string         `bson:"action"        json:"action"`
	CreatedAt    time.Time      `bson:"created_at"    json:"createdAt"`
}

func (h *History) GetID() bson.ObjectID { return h.ID }
func (h *History) Coll() string         { return "history" }

func CreateHistory(userID bson.ObjectID, m db.MongoModel, oldData map[string]any, action string) error {
	history := &History{
		UserID:       userID,
		CollectionID: m.GetID(),
		Collection:   m.Coll(),
		OldData:      oldData,
		Action:       action,
		CreatedAt:    time.Now(),
	}
	if _, e := db.Mongo.Collection(history.Coll()).InsertOne(context.TODO(), history); e != nil {
		return err.Mongo(e)
	}
	return nil
}

// ================================================================
// Tokens
// ================================================================

type Token struct {
	ID        bson.ObjectID `bson:"_id"        json:"id"`
	UserID    string        `bson:"user_id"    json:"userId"`
	Token     string        `bson:"token"      json:"token"`
	Action    string        `bson:"action"     json:"action"`
	CreatedAt time.Time     `bson:"created_at" json:"createdAt"`
	ExpiresAt time.Time     `bson:"expires_at" json:"expiresAt"`
}

func (t *Token) GetID() bson.ObjectID { return t.ID }
func (t *Token) Coll() string         { return "tokens" }

// crea un token nuevo en la bd con un expiresAt de 24 horas
func TokenCreate(user_id string, action string) (*Token, error) {
	t, e := auth.GenerateHexToken()
	if e != nil {
		return nil, e
	}
	tk := &Token{
		UserID:    user_id,
		Token:     t,
		Action:    action,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	result, e := db.Mongo.Collection(tk.Coll()).InsertOne(context.TODO(), tk)
	if e != nil {
		return nil, err.Mongo(e)
	}
	tk.ID = result.InsertedID.(bson.ObjectID)
	return tk, nil
}
