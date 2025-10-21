package schema

import (
	"donbarrigon/new/internal/utils/db"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CreateUsersCollection struct{}

func (u CreateUsersCollection) Name() string {
	return "users"
}

func (u CreateUsersCollection) Up() {
	db.CreateCollection(u.Name(), func(col *db.CollectionBuilder) {
		col.Index("unique", "email")
		col.Index("unique", "nickname")
		col.CreateIndex(mongo.IndexModel{
			Keys: bson.D{{Key: "deleted_at", Value: 1}},
			Options: options.Index().SetPartialFilterExpression(
				bson.D{{Key: "deleted_at", Value: bson.D{{Key: "$exists", Value: false}}}},
			),
		})
	})
}

func (u CreateUsersCollection) Down() {
	db.DropCollection(u.Name())
}
