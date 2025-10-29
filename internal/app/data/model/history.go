package model

import (
	"context"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/err"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ================================================================
// 	                        ACTIONS
// ================================================================

const (
	CREATE_ACTION       = "create"
	UPDATE_ACTION       = "update"
	DELETE_ACTION       = "delete"
	RESTORE_ACTION      = "restore"
	FORCE_DELETE_ACTION = "force-delete"
)

// ================================================================
// 	                      HISTORY MODEL
// ================================================================

type History struct {
	ID           bson.ObjectID `bson:"_id"               json:"id"`
	UserID       bson.ObjectID `bson:"user_id"           json:"userId"`
	CollectionID bson.ObjectID `bson:"collection_id"     json:"collectionId"`
	Collection   string        `bson:"collection"        json:"collection"`
	Changes      *Changes      `bson:"changes,omitempty" json:"changes,omitempty"`
	Action       string        `bson:"action"            json:"action"`
	CreatedAt    time.Time     `bson:"created_at"        json:"createdAt"`
}

// ================================================================
//                FUNCIONES PARA LA INTERFAZ
// ================================================================

func (h *History) GetID() bson.ObjectID { return h.ID }
func (h *History) Coll() string         { return "history" }

// ================================================================
// 	                 FUNCION AUXILIAR
// ================================================================

func CreateHistory(action string, userID bson.ObjectID, m db.MongoModel, changes *Changes) error {
	history := &History{
		UserID:       userID,
		CollectionID: m.GetID(),
		Collection:   m.Coll(),
		Changes:      changes,
		Action:       action,
		CreatedAt:    time.Now(),
	}
	if _, e := db.Mongo.Collection(history.Coll()).InsertOne(context.TODO(), history); e != nil {
		return err.Mongo(e)
	}
	return nil
}

// ================================================================
// 	                        CHANGES
// ================================================================

type Changes struct {
	Old map[string]any
	New map[string]any
}

func NewChanges() *Changes {
	return &Changes{
		Old: map[string]any{},
		New: map[string]any{},
	}
}
