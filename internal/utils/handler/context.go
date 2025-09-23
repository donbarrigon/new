package handler

import (
	"bytes"
	"donbarrigon/new/internal/utils/err"
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
	Can(permissionName ...string) err.Error
	HasRole(roleName ...string) err.Error
}

type Validator interface {
	PrepareForValidation(c *HttpContext) err.Error
}

type MessageResource struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type HttpContext struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Params  map[string]string
	Auth    AuthInterface
}

func NewHttpContext(w http.ResponseWriter, r *http.Request) *HttpContext {
	return &HttpContext{
		Writer:  w,
		Request: r,
	}
}

func (c *HttpContext) Lang() string {
	return c.Request.Header.Get("Accept-Language")
}

func (c *HttpContext) GetBody(request any) err.Error {
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

func (c *HttpContext) ValidateBody(req Validator) err.Error {

	if e := c.GetBody(req); e != nil {
		return e
	}

	errPFV := req.PrepareForValidation(c)
	if errPFV != nil {
		errPFV = errPFV.Errors()
	}
	if err := Validate(c, req); err != nil {
		if errPFV != nil {
			errMap, phMap := errPFV.GetMap()
			for key, valor := range errMap {
				for i, msg := range valor {
					err.Append(&FieldError{
						FieldName:    key,
						Message:      msg,
						Placeholders: phMap[key][i],
					})
				}
			}
		}
		return err.Errors()
	}
	return errPFV
}

func (c *HttpContext) GetParam(param string, defaultValue string) string {
	if value := c.Params[param]; value != "" {
		return value
	}
	return defaultValue
}

func (c *HttpContext) GetInput(param string) string {
	return c.Request.URL.Query().Get(param)
}

func (c *HttpContext) ResponseJSON(status int, data any) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)

	if err := json.NewEncoder(c.Writer).Encode(data); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.WriteHeader(500)
		c.Writer.Write([]byte(Translate(c.Lang(), `{"message": "Error", "error": "Could not encode the response"}`)))
	}
}

func (c *HttpContext) ResponseError(err Error) {
	err.Translate(c.Lang())
	c.ResponseJSON(err.GetStatus(), err)
}

func (c *HttpContext) ResponseNotFound() {
	c.ResponseError(Errors.NotFoundf("The resource [{method}:{path}] does not exist",
		Entry{"method", c.Request.Method},
		Entry{"path", c.Request.URL.Path},
	))
}

func (c *HttpContext) ResponseMessage(code int, data any, message string, ph ...Entry) {
	c.ResponseJSON(code, &MessageResource{
		Message: Translate(c.Lang(), message, ph...),
		Data:    data,
	})
}

func (c *HttpContext) ResponseOk(data any) {
	c.ResponseJSON(http.StatusOK, data)
}

func (c *HttpContext) ResponseCreated(data any) {
	c.ResponseJSON(http.StatusCreated, data)
}

func (c *HttpContext) ResponseNoContent() {
	c.Writer.WriteHeader(http.StatusNoContent)
}

func (c *HttpContext) ResponseCSV(fileName string, data any, comma ...rune) {
	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Slice {
		err := &Err{
			Status:  http.StatusInternalServerError,
			Message: "Error writing CSV",
			Err:     "Data is not a slice of structs",
		}
		c.ResponseError(err)
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
		err := Errors.NoDocumentsf("No data available")
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
