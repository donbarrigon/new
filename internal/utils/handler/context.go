package handler

import (
	"bytes"
	"donbarrigon/new/internal/utils/auth"
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/lang"
	"donbarrigon/new/internal/utils/str"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	handler *Handler
	Auth    *auth.Session
}

func NewContext(w http.ResponseWriter, r *http.Request, h *Handler) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		handler: h,
	}
}

func (c *Context) Lang() string {
	return c.Request.Header.Get("Accept-Language")
}

func (c *Context) GetBody(request any) error {
	decoder := json.NewDecoder(c.Request.Body)
	if e := decoder.Decode(request); e != nil {
		return err.New(
			http.StatusBadRequest,
			"El cuerpo de la solicitud no es válido",
			e.Error(),
		)
	}
	defer c.Request.Body.Close()
	return nil
}

func (c *Context) Get(param string) string {
	return c.Request.URL.Query().Get(param)
}

func (c *Context) ResponseJson(status int, data any) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)

	if err := json.NewEncoder(c.Writer).Encode(data); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.WriteHeader(500)
		c.Writer.Write([]byte(lang.T(c.Lang(), `{"message": "Error", "error": "Could not encode the response"}`, nil)))
	}
}

func (c *Context) Response(status int, data any) {
	c.Writer.Header().Set("Content-Type", "application/msgpack")
	c.Writer.WriteHeader(status)

	if err := msgpack.NewEncoder(c.Writer).Encode(data); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(`{"message": "Error", "error": "Could not encode the response"}`))
	}
}

func (c *Context) Responsef(status int, data any) {
	acept := c.Request.Header.Get("accept")
	if acept == "application/json" {
		c.ResponseJson(status, data)
		return
	}
	c.Response(status, data)
}

func (c *Context) ResponseError(e error) {

	var he *err.HttpError
	if errors.As(e, &he) {
		he.Message = lang.T(c.Lang(), he.Message, nil)
		c.Responsef(he.Status, he)
		return
	}
	var ve *err.ValidationError
	if errors.As(e, &ve) {
		her := ve.Herror(c.Lang())
		c.Responsef(her.Status, her)
		return
	}
	her := err.Internal(e)
	c.Responsef(her.Status, her)
}

func (c *Context) ResponseNotFound() {
	c.ResponseError(err.NotFound(lang.T(c.Lang(), "El recurso [:method :path] No existe", str.Placeholder{
		{Key: "method", Value: c.Request.Method},
		{Key: "path", Value: c.Request.URL.Path},
	})))
}

func (c *Context) ResponseOk(data any) {
	acept := c.Request.Header.Get("accept")
	if acept == "application/json" {
		c.ResponseJson(http.StatusOK, data)
		return
	}
	c.Response(http.StatusOK, data)
}

func (c *Context) ResponseCreated(data any) {
	acept := c.Request.Header.Get("accept")
	if acept == "application/json" {
		c.ResponseJson(http.StatusCreated, data)
		return
	}
	c.Response(http.StatusCreated, data)
}

func (c *Context) ResponseNoContent() {
	c.Writer.WriteHeader(http.StatusNoContent)
}

func (c *Context) ResponseCSV(fileName string, data any, comma ...rune) {
	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Slice {
		e := &err.HttpError{
			Status:  http.StatusInternalServerError,
			Message: "Error writing CSV",
			Err:     "Data is not a slice of structs",
		}
		c.ResponseError(e)
		return
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	del := ';'
	if len(comma) > 0 {
		del = comma[0]
	}
	writer.Comma = del

	if val.Len() == 0 {
		err := err.NotFound("No data available")
		c.ResponseError(err)
		return
	}

	first := val.Index(0)
	elemType := first.Type()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	var headers []string
	var fields []int

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		tag = strings.Split(tag, ",")[0]
		headers = append(headers, tag)
		fields = append(fields, i)
	}
	writer.Write(headers)

	for i := 0; i < val.Len(); i++ {
		var record []string
		elem := val.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		for _, j := range fields {
			fieldVal := elem.Field(j)

			if fieldVal.Type() == reflect.TypeOf(bson.ObjectID{}) {
				objID := fieldVal.Interface().(bson.ObjectID)
				record = append(record, objID.Hex()) // sin comillas manuales
				continue
			}

			switch fieldVal.Kind() {
			case reflect.String:
				record = append(record, fieldVal.String())
			case reflect.Int, reflect.Int64:
				record = append(record, fmt.Sprintf("%d", fieldVal.Int()))
			case reflect.Float64:
				record = append(record, fmt.Sprintf("%f", fieldVal.Float()))
			case reflect.Bool:
				record = append(record, fmt.Sprintf("%t", fieldVal.Bool()))
			case reflect.Struct:
				if t, ok := fieldVal.Interface().(time.Time); ok {
					record = append(record, t.Format(time.RFC3339))
				} else {
					jsonVal, _ := json.Marshal(fieldVal.Interface())
					record = append(record, string(jsonVal))
				}
			case reflect.Slice, reflect.Map, reflect.Array:
				jsonVal, _ := json.Marshal(fieldVal.Interface())
				record = append(record, string(jsonVal))
			default:
				record = append(record, fmt.Sprintf("%v", fieldVal.Interface()))
			}
		}
		writer.Write(record)
	}
	writer.Flush()

	c.Writer.Header().Set("Content-Type", "text/csv")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename="+fileName+".csv")
	c.Writer.Write(buffer.Bytes())
}
