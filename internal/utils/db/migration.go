package db

import (
	"context"
	"donbarrigon/new/internal/utils/fm"
	"donbarrigon/new/internal/utils/logs"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Migration interface {
	name() string
	Up() func()
	Down() func()
}

type CollectionBuilder struct {
	Name string
}

func CreateCollection(name string, callback func(col *CollectionBuilder)) {
	col := &CollectionBuilder{Name: name}
	callback(col)
}

func AlterCollection(name string, callback func(col *CollectionBuilder)) {
	col := &CollectionBuilder{Name: name}
	callback(col)
}

func (c *CollectionBuilder) Index(fields ...string) {
	keys := bson.D{}
	fm := "[ "
	for _, field := range fields {
		fs := strings.Split(field, ":")
		f := fs[0]
		s := 1
		if len(fs) == 2 {
			s, _ = strconv.Atoi(fs[1])
		}
		keys = append(keys, bson.E{Key: f, Value: s})
		fm += f + " "
	}
	fm += "]"

	name, e := DB.Collection(c.Name).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: keys,
	})
	if e != nil {
		logs.Error("Error al crear el indice %s %s %s", c.Name, fm, e.Error())
		panic(e.Error())
	}
	logs.Info("Indice creado %s %s", c.Name, name)
}

func (c *CollectionBuilder) UniqueIndex(fields ...string) {

	keys := bson.D{}
	fm := "[ "
	for _, field := range fields {
		fs := strings.Split(field, ":")
		f := fs[0]
		s := 1
		if len(fs) == 2 {
			s, _ = strconv.Atoi(fs[1])
		}
		keys = append(keys, bson.E{Key: f, Value: s})
		fm += f + " "
	}
	fm += "]"

	name, e := DB.Collection(c.Name).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true),
	})
	if e != nil {
		logs.Error("Error al crear el indice unico %s %s %s", c.Name, fm, e.Error())
		panic(e.Error())
	}
	logs.Info("Indice unico creado %s %s", c.Name, name)
}

func (c *CollectionBuilder) TextIndex(fields ...string) {
	keys := bson.D{}
	fm := "[ "
	for _, field := range fields {
		keys = append(keys, bson.E{Key: field, Value: "text"})
		fm += field + " "
	}
	fm += "]"

	name, e := DB.Collection(c.Name).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: keys,
	})
	if e != nil {
		logs.Error("Error al crear el indice de texto %s %s %s", c.Name, fm, e.Error())
		panic(e.Error())
	}
	logs.Info("Indice de texto creado %s %s", c.Name, name)
}

func (c *CollectionBuilder) HashedIndex(field string) {
	name, e := DB.Collection(c.Name).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{Key: field, Value: "hashed"}},
	})
	if e != nil {
		logs.Error("Error al crear el indice hashed %s %s %s", c.Name, field, e.Error())
		panic(e.Error())
	}
	logs.Info("Indice hashed creado %s %s", c.Name, name)
}

func (c *CollectionBuilder) CreateOneIndex(model mongo.IndexModel, opts ...options.Lister[options.CreateIndexesOptions]) {
	name, e := DB.Collection(c.Name).Indexes().CreateOne(context.TODO(), model, opts...)
	if e != nil {
		logs.Error("Error al crear el indice :collection :error", fm.Placeholder{
			"collection": c.Name,
			"error":      e.Error(),
		})
	}
	logs.Info("Indice creado :collection :name", fm.Placeholder{
		"collection": c.Name,
		"name":       name,
	})
}

func (c *CollectionBuilder) DropIndex(indexName string) {

	e := DB.Collection(c.Name).Indexes().DropOne(context.TODO(), indexName)
	if e != nil {
		logs.Error("Error al eliminar el indice :collection :name :error ", fm.Placeholder{
			"collection": c.Name,
			"name":       indexName,
			"error":      e.Error(),
		})
		panic(e.Error())
	}
	logs.Info("Indice eliminado :collection :name", fm.Placeholder{
		"collection": c.Name,
		"name":       indexName,
	})
}

func (c *CollectionBuilder) DropAllIndexes() {

	e := DB.Collection(c.Name).Indexes().DropAll(context.TODO())
	if e != nil {
		logs.Error("Error al eliminar todos los indices :collection :error ", fm.Placeholder{
			"collection": c.Name,
			"error":      e.Error(),
		})
		panic(e.Error())
	}
	logs.Info("Todos los indices fueron eliminados :collection", fm.Placeholder{
		"collection": c.Name,
	})
}

func DropCollection(collection string) {

	e := DB.Collection(collection).Drop(context.TODO())
	if e != nil {
		logs.Error("Failed to drop collection :collection :error ", fm.Placeholder{
			"collection": collection,
			"error":      e.Error(),
		})
		panic(e.Error())
	}
	logs.Info("Dropped collection :collection", fm.Placeholder{
		"collection": collection,
	})
}
