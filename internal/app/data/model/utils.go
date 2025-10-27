package model

import (
	"donbarrigon/new/internal/utils/auth"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/err"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ================================================================
// trash
// ================================================================
type Trash struct {
	ID           bson.ObjectID `bson:"_id"           json:"id"`
	UserID       bson.ObjectID `bson:"user_id"       json:"userId"`
	Collection   string        `bson:"collection"    json:"collection"`
	CollectionID bson.ObjectID `bson:"collection_id" json:"collectionId"`
	Data         any           `bson:"data"          json:"data"`
	DeletedAt    time.Time     `bson:"deleted_at"    json:"deletedAt"`
	db.Odm       `bson:"-" json:"-"`
}

func (t *Trash) CollectionName() string { return "trash" }
func (t *Trash) GetID() bson.ObjectID   { return t.ID }
func (t *Trash) SetID(id bson.ObjectID) { t.ID = id }

func MoveToTrash(userID bson.ObjectID, m db.OdmModel) error {
	trash := &Trash{
		UserID:       userID,
		Collection:   m.CollectionName(),
		CollectionID: m.GetID(),
		Data:         m,
	}
	trash.Odm.Model = trash
	if e := trash.Create(); e != nil {
		return e
	}
	if e := m.Delete(); e != nil {
		return e
	}
	return nil
}

// ================================================================
// History
// ================================================================

type History struct {
	ID         bson.ObjectID  `bson:"_id"        json:"id"`
	UserID     bson.ObjectID  `bson:"user_id"    json:"userId"`
	Collection string         `bson:"collection" json:"collection"`
	OldData    map[string]any `bson:"old_data"   json:"oldData"`
	NewData    map[string]any `bson:"new_data"   json:"newData"`
	Action     string         `bson:"action"     json:"action"`
	CreatedAt  time.Time      `bson:"created_at" json:"createdAt"`
	db.Odm     `bson:"-" json:"-"`
}

func (h *History) CollectionName() string { return "history" }
func (h *History) GetID() bson.ObjectID   { return h.ID }
func (h *History) SetID(id bson.ObjectID) { h.ID = id }
func (h *History) BeforeCreate() error {
	h.CreatedAt = time.Now()
	return nil
}
func (h *History) BeforeUpdate() error {
	return err.New(err.INTERNAL, "Estan intentando modificar el historial", errors.New("Estan intentando modificar el historial"))
}

func NewHistory(userID bson.ObjectID, m db.OdmModel, action string) error {
	history := &History{
		UserID:     userID,
		Collection: m.CollectionName(),
		OldData:    m.GetOriginal(),
		NewData:    m.GetDirty(),
		Action:     action,
	}
	history.Odm.Model = history
	if e := history.Create(); e != nil {
		return e
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
	db.Odm    `bson:"-" json:"-"`
}

// crea un token nuevo en la bd con un expiresAt de 24 horas
func NewToken(user_id string, action string) (*Token, error) {
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
	tk.Odm.Model = tk
	if e := tk.Create(); e != nil {
		return nil, e
	}
	return tk, nil
}

func (t *Token) CollectionName() string { return "tokens" }
func (t *Token) GetID() bson.ObjectID   { return t.ID }
func (t *Token) SetID(id bson.ObjectID) { t.ID = id }
