package db

import "go.mongodb.org/mongo-driver/v2/mongo"

func Col(col string) *mongo.Collection {
	return DB.Collection(col)
}
