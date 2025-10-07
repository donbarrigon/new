package handler

import (
	"bytes"
	"donbarrigon/new/internal/utils/auth"
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/fm"
	"donbarrigon/new/internal/utils/lang"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthInterface interface {
	GetID() bson.ObjectID
	GetUserID() bson.ObjectID
	Can(permissionName string) err.Error
	HasRole(roleName string) err.Error
}

type Validator interface {
	PrepareForValidation(hc *Context) err.Error
}

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	handler *Handler
	Session *auth.Session
}

func NewContext(w http.ResponseWriter, r *http.Request, h *Handler) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		handler: h,
	}
}

func (hc *Context) Lang() string {
	return hc.Request.Header.Get("Accept-Language")
}

func (hc *Context) GetBody(request any) err.Error {
	decoder := json.NewDecoder(hc.Request.Body)
	if e := decoder.Decode(request); e != nil {
		return err.New(
			http.StatusBadRequest,
			"El cuerpo de la solicitud no es válido",
			e.Error(),
		)
	}
	defer hc.Request.Body.Close()
	return nil
}

func (hc *Context) Get(param string, defaultValue string) string {
	return hc.Request.URL.Query().Get(param)
}

func (hc *Context) ResponseJSON(status int, data any) {
	hc.Writer.Header().Set("Content-Type", "application/json")
	hc.Writer.WriteHeader(status)

	if err := json.NewEncoder(hc.Writer).Encode(data); err != nil {
		hc.Writer.WriteHeader(http.StatusInternalServerError)
		hc.Writer.WriteHeader(500)
		hc.Writer.Write([]byte(lang.T(hc.Lang(), `{"message": "Error", "error": "Could not encode the response"}`, nil)))
	}
}

func (hc *Context) ResponseError(e err.Error) {
	er := e.Errors(hc.Lang())
	hc.ResponseJSON(er.Status, er)
}

func (hc *Context) ResponseNotFound() {
	hc.ResponseError(err.NotFound(lang.T(hc.Lang(), "The resource [:method :path] does not exist", fm.Placeholder{"method": hc.Request.Method, "path": hc.Request.URL.Path})))
}

func (hc *Context) ResponseOk(data any) {
	hc.ResponseJSON(http.StatusOK, data)
}

func (hc *Context) ResponseCreated(data any) {
	hc.ResponseJSON(http.StatusCreated, data)
}

func (hc *Context) ResponseNoContent() {
	hc.Writer.WriteHeader(http.StatusNoContent)
}

func (hc *Context) ResponseCSV(fileName string, data any, comma ...rune) {
	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Slice {
		err := &err.HttpError{
			Status:  http.StatusInternalServerError,
			Message: "Error writing CSV",
			Err:     "Data is not a slice of structs",
		}
		hc.ResponseError(err)
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
		hc.ResponseError(err)
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

	hc.Writer.Header().Set("Content-Type", "text/csv")
	hc.Writer.Header().Set("Content-Disposition", "attachment;filename="+fileName+".csv")
	hc.Writer.Write(buffer.Bytes())
}
