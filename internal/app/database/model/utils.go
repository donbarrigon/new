package model

import (
	"donbarrigon/new/internal/utils/auth"
	"donbarrigon/new/internal/utils/db"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ================================================================
// trash
// ================================================================
type Trash struct {
	ID         bson.ObjectID `bson:"_id" json:"id"`
	UserID     bson.ObjectID `bson:"user_id" json:"userId"`
	Collection string        `bson:"collection" json:"collection"`
	Data       any           `bson:"data" json:"data"`
	DeletedAt  time.Time     `bson:"deleted_at" json:"deletedAt"`
	db.Odm     `bson:"-" json:"-"`
}

func (t *Trash) CollectionName() string { return "trash" }
func (t *Trash) GetID() bson.ObjectID   { return t.ID }
func (t *Trash) SetID(id bson.ObjectID) { t.ID = id }

func MoveToTrash(user_id bson.ObjectID, col db.Model) error {
	trash := &Trash{
		UserID:     user_id,
		Collection: col.CollectionName(),
		Data:       col,
		DeletedAt:  time.Now(),
	}
	trash.Odm.Model = trash
	if e := trash.Create(); e != nil {
		return e
	}
	if e := col.Delete(); e != nil {
		return e
	}
	return nil
}

// ================================================================
// Tokens
// ================================================================

type Token struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	UserID    string        `bson:"user_id" json:"userId"`
	Token     string        `bson:"token" json:"token"`
	Action    string        `bson:"action" json:"action"`
	CreatedAt time.Time     `bson:"created_at" json:"createdAt"`
	ExpiresAt time.Time     `bson:"expires_at" json:"expiresAt"`
	db.Odm    `bson:"-" json:"-"`
}

// crea un token nuevo en la bd con un expiresAt de 24 horas
func NewToken(user_id string, action string) (*Token, error) {
	t, e := auth.GenerateToken()
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
	tk.Odm.Model = tk
	if e := tk.Create(); e != nil {
		return nil, e
	}
	return tk, nil
}

func (t *Token) CollectionName() string { return "tokens" }
func (t *Token) GetID() bson.ObjectID   { return t.ID }
func (t *Token) SetID(id bson.ObjectID) { t.ID = id }
