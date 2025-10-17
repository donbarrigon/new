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
	BeforeCreate() error
	BeforeUpdate() error
	BeforeDelete() error
	AfterCreate() error
	AfterUpdate() error
	AfterDelete() error

	Create() error // esto lo agrega el odm
	Delete() error // esto lo agrega el odm
}

type Collection []Model

type Odm struct {
	Model Model `bson:"-" json:"-"`
}

func (o *Odm) BeforeCreate() error { return nil }
func (o *Odm) BeforeUpdate() error { return nil }
func (o *Odm) BeforeDelete() error { return nil }
func (o *Odm) AfterCreate() error  { return nil }
func (o *Odm) AfterUpdate() error  { return nil }
func (o *Odm) AfterDelete() error  { return nil }

func (o *Odm) FindByHexID(id string) error {

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

func (o *Odm) FindByID(id bson.ObjectID) error {
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	if e := DB.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter).Decode(o.Model); e != nil {
		return err.Mongo(e)
	}
	return nil
}

func (o *Odm) First(field string, value any) error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	if e := DB.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter).Decode(o.Model); e != nil {
		return err.Mongo(e)
	}
	return nil
}

func (o *Odm) FindOne(filter bson.D, opts ...options.Lister[options.FindOneOptions]) error {
	if e := DB.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter, opts...).Decode(o.Model); e != nil {
		return err.Mongo(e)
	}
	return nil
}

func (o *Odm) Find(result any, filter bson.D, opts ...options.Lister[options.FindOptions]) error {
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

// busqueda eq
func (o *Odm) FindByField(result any, field string, value any, opts ...options.Lister[options.FindOptions]) error {
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

func (o *Odm) Aggregate(result any, pipeline mongo.Pipeline) error {
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

func (o *Odm) AggregateOne(pipeline mongo.Pipeline) error {
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
		return err.New(err.NOT_FOUND, "El documento no existe", "!cursor.Next(ctx)")
	}
	return nil
}

func (o *Odm) Create() error {
	if e := o.Model.BeforeCreate(); e != nil {
		return e
	}
	result, e := DB.Collection(o.Model.CollectionName()).InsertOne(context.TODO(), o.Model)
	if e != nil {
		return err.Mongo(e)
	}
	o.Model.SetID(result.InsertedID.(bson.ObjectID))
	return o.Model.AfterCreate()
}

func (o *Odm) CreateBy(validator any) error {
	if e := Fill(o.Model, validator); e != nil {
		return e
	}
	return o.Create()
}

// usela solo si tienes pereza.
func (o *Odm) CreateMany(data any) error {

	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Slice {
		return err.New(err.INTERNAL, "No se puede guardar la coleccion de datos", "CreateMany solo acepta slices de Model")
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
	he := []error{}
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		elem.(Model).SetID(result.InsertedIDs[i].(bson.ObjectID))
		if e := elem.(Model).AfterCreate(); e != nil {
			he = append(he, e)
		}
	}
	if len(he) > 0 {
		return err.New(err.INTERNAL, "Algo salio mal al guardar la coleccion de datos", he)
	}
	return nil
}

func (o *Odm) Update() error {
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
		return err.New(err.NOT_FOUND, "El documento a modificar no existe", "!result.MatchedCount == 0")
	}

	if result.ModifiedCount == 0 {
		return err.New(err.CONFLICT, "No se aplicaron cambios al guardar el documento", "!result.ModifiedCount == 0")
	}
	return o.Model.AfterUpdate()
}

func (o *Odm) UpdateBy(validator any) (map[string]any, map[string]any, error) {
	original, dirty, e := Filld(o.Model, validator)
	if e != nil {
		return original, dirty, e
	}
	return original, dirty, o.Update()
}

// OjO no usa el hook BeforeUpdate

func (o *Odm) Delete() error {
	if e := o.Model.BeforeDelete(); e != nil {
		return e
	}

	filter := bson.D{bson.E{Key: "_id", Value: o.Model.GetID()}}

	result, e := DB.Collection(o.Model.CollectionName()).DeleteOne(context.TODO(), filter)
	if e != nil {
		return err.Mongo(e)
	}
	if result.DeletedCount == 0 {
		return err.New(err.CONFLICT, "No se elimino el documento", "!result.DeletedCount == 0")
	}
	return o.Model.AfterDelete()
}
