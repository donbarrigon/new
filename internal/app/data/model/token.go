package model

import (
	"context"
	"donbarrigon/new/internal/utils/auth"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/err"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Token struct {
	ID        bson.ObjectID     `bson:"_id"        json:"id"`
	UserID    bson.ObjectID     `bson:"user_id"    json:"userId"`
	Token     string            `bson:"token"      json:"token"`
	Action    string            `bson:"action"     json:"action"`
	Metadata  map[string]string `bson:"metadata"   json:"metadata"`
	CreatedAt time.Time         `bson:"created_at" json:"createdAt"`
	ExpiresAt time.Time         `bson:"expires_at" json:"expiresAt"`
}

// ================================================================
//                 FUNCIONES PARA LA INTERFAZ
// ================================================================

func (t *Token) GetID() bson.ObjectID { return t.ID }
func (t *Token) Coll() string         { return "tokens" }

// ================================================================
//                     FUNCIONES CRUD
// ================================================================

// crea un token nuevo en la bd con un expiresAt de 24 horas
func TokenCreate(user_id bson.ObjectID, action string, metadata map[string]string) (*Token, error) {
	t, e := auth.GenerateHexToken()
	if e != nil {
		return nil, e
	}
	tk := &Token{
		UserID:    user_id,
		Token:     t,
		Action:    action,
		Metadata:  metadata,
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
