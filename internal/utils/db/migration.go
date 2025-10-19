package db

import (
	"context"
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/fm"
	"donbarrigon/new/internal/utils/logs"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Migration interface {
	Name() string
	Up()
	Down()
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

// Index crea un indice en la coleccion
// tipe: i, index, u, unique, t, text, h, hashed, s, sparse
// fields: los campos a indexar por defecto el orden es 1
// fields: el orden decendente se puede colocar con fieldName:-1
// Forma de uso
// col.Index("u", "email") // agrega un indice unico a email
// col.Index("i", "created_at:-1") // agrega un indice en orden desc
// col.Index("i", "name:1", "last_name:-1") // agrega un indice compuesto a name y last_name
// col.Index("t", "city", "state", "country") // agrega un indice de texto a city, state y country
// col.Index("h", "consecutive") // agrega un indice hashed el hashed solo recibe un campo
// col.Index("s", "deleted_at") // agrega un indice sparce el sparce solo recibe un campo
func (c *CollectionBuilder) Index(tipe string, fields ...string) {
	if len(fields) > 0 {
		logs.Error("No colocaste campos para crear el indice %s %s [%s]", c.Name, tipe, strings.Join(fields, ","))
		panic(err.New(err.INTERNAL, "No colocaste campos para crear el indice", nil))
	}
	switch tipe {
	case "i", "index":
		c.indexOne(fields...)
	case "u", "unique":
		c.indexUnique(fields...)
	case "t", "text":
		c.indexText(fields...)
	case "h", "hashed":
		c.indexHashed(fields[0])
	case "s", "sparce":
		c.indexSparce(fields[0])
	default:
		logs.Error("No se reconoce el tipo de indice %s %s [%s]", c.Name, tipe, strings.Join(fields, ","))
		panic(err.New(err.INTERNAL, "No se reconoce el indice", nil))
	}
}
func (c *CollectionBuilder) indexOne(fields ...string) {
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

func (c *CollectionBuilder) indexUnique(fields ...string) {

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

func (c *CollectionBuilder) indexSparce(fields ...string) {

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
		Options: options.Index().SetSparse(true),
	})
	if e != nil {
		logs.Error("Error al crear el indice sparce %s %s %s", c.Name, fm, e.Error())
		panic(e.Error())
	}
	logs.Info("Indice sparce creado %s %s", c.Name, name)
}

func (c *CollectionBuilder) indexText(fields ...string) {
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

func (c *CollectionBuilder) indexHashed(field string) {
	name, e := DB.Collection(c.Name).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{Key: field, Value: "hashed"}},
	})
	if e != nil {
		logs.Error("Error al crear el indice hashed %s %s %s", c.Name, field, e.Error())
		panic(e.Error())
	}
	logs.Info("Indice hashed creado %s %s", c.Name, name)
}

func (c *CollectionBuilder) CreateIndex(model mongo.IndexModel, opts ...options.Lister[options.CreateIndexesOptions]) {
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
