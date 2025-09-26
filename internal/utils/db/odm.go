package db

import (
	"context"
	"donbarrigon/new/internal/utils/err"
	"reflect"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Model interface {
	CollectionName() string
	GetID() bson.ObjectID
	SetID(id bson.ObjectID)
	BeforeCreate() err.Error
	BeforeUpdate() err.Error

	Create() err.Error // esto lo agrega el odm
	Delete() err.Error // esto lo agrega el odm
}

type Collection []Model

type Odm struct {
	Model Model `bson:"-" json:"-"`
}

var DBClient *mongo.Client
var DB *mongo.Database

func (o *Odm) FindByHexID(id string) err.Error {

	objectId, e := bson.ObjectIDFromHex(id)
	if e != nil {
		return err.HexID(e.Error())
	}
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	if e := DB.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter).Decode(o.Model); e != nil {
		return err.Mongo(e)
	}
	return nil
}

func (o *Odm) FindByID(id bson.ObjectID) err.Error {
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	if e := DB.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter).Decode(o.Model); e != nil {
		return err.Mongo(e)
	}
	return nil
}

func (o *Odm) First(field string, value any) err.Error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	if e := DB.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter).Decode(o.Model); e != nil {
		return err.Mongo(e)
	}
	return nil
}

func (o *Odm) FindOne(filter bson.D, opts ...options.Lister[options.FindOneOptions]) err.Error {
	if e := DB.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter, opts...).Decode(o.Model); e != nil {
		return err.Mongo(e)
	}
	return nil
}

func (o *Odm) Find(result any, filter bson.D, opts ...options.Lister[options.FindOptions]) err.Error {
	ctx := context.TODO()
	cursor, e := DB.Collection(o.Model.CollectionName()).Find(ctx, filter, opts...)
	if e != nil {
		return err.Mongo(e)
	}
	if e = cursor.All(ctx, result); e != nil {
		return err.Mongo(e)
	}
	return nil
}

func (o *Odm) FindBy(result any, field string, value any, opts ...options.Lister[options.FindOptions]) err.Error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	ctx := context.TODO()
	cursor, e := DB.Collection(o.Model.CollectionName()).Find(ctx, filter, opts...)
	if e != nil {
		return err.Mongo(e)
	}
	if e = cursor.All(ctx, result); e != nil {
		return err.Mongo(e)
	}
	return nil
}

func (o *Odm) Aggregate(result any, pipeline mongo.Pipeline) err.Error {
	ctx := context.TODO()
	cursor, e := DB.Collection(o.Model.CollectionName()).Aggregate(ctx, pipeline)
	if e != nil {
		return err.Mongo(e)
	}
	if e = cursor.All(ctx, result); e != nil {
		return err.Mongo(e)
	}
	return nil
}

func (o *Odm) AggregateOne(pipeline mongo.Pipeline) err.Error {
	ctx := context.TODO()
	cursor, e := DB.Collection(o.Model.CollectionName()).Aggregate(ctx, pipeline)
	if e != nil {
		return err.Mongo(e)
	}
	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		if e := cursor.Decode(o.Model); e != nil {
			return err.Mongo(e)
		}
	} else {
		return err.NotFound("document not found")
	}
	return nil
}

func (o *Odm) Create() err.Error {
	if e := o.Model.BeforeCreate(); e != nil {
		return e
	}
	result, e := DB.Collection(o.Model.CollectionName()).InsertOne(context.TODO(), o.Model)
	if e != nil {
		return err.Mongo(e)
	}
	o.Model.SetID(result.InsertedID.(bson.ObjectID))

	return nil
}

func (o *Odm) CreateBy(validator any) err.Error {
	if _, _, e := Fill(o.Model, validator); e != nil {
		return e
	}
	return o.Create()
}

func (o *Odm) CreateMany(data any) err.Error {

	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Slice {
		return err.Internal("Create many required a slice")
	}
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		if e := elem.(Model).BeforeCreate(); e != nil {
			return e
		}
	}
	collection := DB.Collection(o.Model.CollectionName())
	result, e := collection.InsertMany(context.TODO(), data)
	if e != nil {
		return err.Mongo(e)
	}
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		elem.(Model).SetID(result.InsertedIDs[i].(bson.ObjectID))
	}
	return nil
}

func (o *Odm) Update() err.Error {
	if e := o.Model.BeforeUpdate(); e != nil {
		return e
	}
	filter := bson.D{bson.E{Key: "_id", Value: o.Model.GetID()}}
	update := bson.D{bson.E{Key: "$set", Value: o.Model}}

	result, e := DB.Collection(o.Model.CollectionName()).UpdateOne(context.TODO(), filter, update)
	if e != nil {
		return err.Mongo(e)
	}
	if result.MatchedCount == 0 {
		return err.NotFound("document not found for update")
	}

	if result.ModifiedCount == 0 {
		return err.Conflict("no changes were applied to the document")
	}
	return nil
}

func (o *Odm) UpdateBy(validator any) (map[string]any, map[string]any, err.Error) {
	original, dirty, e := Fill(o.Model, validator)
	if e != nil {
		return original, dirty, e
	}
	return original, dirty, o.Update()
}

// OjO no usa el hook BeforeUpdate

func (o *Odm) Delete() err.Error {
	filter := bson.D{bson.E{Key: "_id", Value: o.Model.GetID()}}

	result, e := DB.Collection(o.Model.CollectionName()).DeleteOne(context.TODO(), filter)
	if e != nil {
		return err.Mongo(e)
	}
	if result.DeletedCount == 0 {
		return err.Conflict("no changes were applied to the document")
	}
	return nil
}
