package model

import (
	"context"
	"donbarrigon/new/internal/utils/auth"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/logs"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	result, e := db.Mongo.Collection(tk.Coll()).InsertOne(context.TODO(), tk)
	if e != nil {
		return nil, err.Mongo(e)
	}
	tk.ID = result.InsertedID.(bson.ObjectID)
	return tk, nil
}

func TokenCreateVerificationCode(userID bson.ObjectID, action string, metadata map[string]string) (*Token, error) {
	t, e := auth.GenerateVerificationCode()
	if e != nil {
		return nil, e
	}
	tk := &Token{
		UserID:    userID,
		Token:     t,
		Action:    action,
		Metadata:  metadata,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	result, e := db.Mongo.Collection(tk.Coll()).InsertOne(context.TODO(), tk)
	if e != nil {
		return nil, err.Mongo(e)
	}
	tk.ID = result.InsertedID.(bson.ObjectID)
	return tk, nil
}

func TokenGet(action string, userID bson.ObjectID, token string) (*Token, error) {
	filter := bson.D{
		{Key: "user_id", Value: userID},
		{Key: "token", Value: token},
		{Key: "action", Value: action},
	}

	tk := &Token{}
	if e := db.Mongo.Collection(tk.Coll()).FindOne(context.TODO(), filter).Decode(tk); e != nil {
		if e == mongo.ErrNoDocuments {
			return nil, err.New(err.UNAUTHORIZED, "El token es incorrecto", nil)
		}
		return nil, err.Mongo(e)
	}

	return tk, nil
}

// busca si existe y no a expitado y lo elimina
// retorna nil si todo salio bien
func TokenExists(action string, userID bson.ObjectID, token string) error {
	filter := bson.D{
		{Key: "user_id", Value: userID},
		{Key: "token", Value: token},
		{Key: "action", Value: action},
	}

	tk := &Token{}
	if e := db.Mongo.Collection(tk.Coll()).FindOne(context.TODO(), filter).Decode(tk); e != nil {
		if e == mongo.ErrNoDocuments {
			return err.New(err.UNAUTHORIZED, "El token es incorrecto", nil)
		}
		return err.Mongo(e)
	}

	if tk.ExpiresAt.Before(time.Now()) {
		if e := tk.Delete(); e != nil {
			logs.Error("No se pudo eliminar el token %s : %s", tk.ID.Hex(), e.Error())
		}
		return err.New(err.UNAUTHORIZED, "El token ha expirado", nil)
	}

	if e := tk.Delete(); e != nil {
		logs.Error("No se pudo eliminar el token %s : %s", tk.ID.Hex(), e.Error())
	}

	return nil
}

func (tk *Token) Delete() error {
	filter := bson.D{{Key: "_id", Value: tk.ID}}
	result, e := db.Mongo.Collection(tk.Coll()).DeleteOne(context.TODO(), filter)
	if e != nil {
		return err.Mongo(e)
	}
	if e := err.MongoDeleteResult(result); e != nil {
		return e
	}
	return nil
}
