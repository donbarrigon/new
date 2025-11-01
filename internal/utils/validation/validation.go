package validation

import (
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/handler"
	"donbarrigon/new/internal/utils/str"
	"reflect"
	"strings"
)

type ValidationFunc func(value reflect.Value, params ...string) (string, str.Placeholder, bool)

type Rules map[string]map[string][]string

type Validator interface {
	Rules() Rules
	PrepareForValidation(c *handler.Context) *err.ValidationError
}

func Body(c *handler.Context, validator Validator) error {
	if e := c.GetBody(validator); e != nil {
		return e
	}

	e := validator.PrepareForValidation(c)
	// if e == nil {
	// 	e = err.NewValidationError()
	// }

	rules := validator.Rules()
	return validate(c, validator, rules, e)
}

func From(c *handler.Context, data any, rules Rules) error {
	e := err.NewValidationError()
	return validate(c, data, rules, e)
}

func validate(c *handler.Context, validator any, rules Rules, e *err.ValidationError) error {
	val := reflect.ValueOf(validator)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return err.Internal("Error de validación: Los datos tienen que ser un *struct")
	}
	typ := val.Type()
	numFields := typ.NumField()
	for i := range numFields {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		tagName := strings.Split(tag, ",")[0]
		if tagName == "-" || tagName == "id" {
			continue
		}
		if tagName == "" {
			tagName = field.Name
		}

		value := val.Field(i)
		if value.Kind() == reflect.Pointer {
			if !value.IsNil() {
				value = value.Elem()
			}
		}

		if value.Kind() == reflect.Struct {
			r := Rules{}
			for key, rule := range rules {
				if k, found := strings.CutPrefix(key, tagName+"."); found {
					r[k] = rule
				}
			}
			e2 := err.NewValidationError()
			if validate(c, value.Interface(), r, e2) != nil {
				for e2k, e2v := range e2.Messages {
					e.Messages[tagName+"."+e2k] = e2v
					e.Placeholders[tagName+"."+e2k] = e2.Placeholders[e2k]
				}
			}
			continue
		}

		validations := rules[tagName]
		if validations == nil {
			continue
		}
		for rule, params := range validations {
			switch rule {
			case "required":
				if msg, ph, hasError := Required(value, params...); hasError {
					e.Append(tagName, msg, ph)
				}
			case "min":
				if msg, ph, hasError := Min(value, params...); hasError {
					e.Append(tagName, msg, ph)
				}
			case "max":
				if msg, ph, hasError := Max(value, params...); hasError {
					e.Append(tagName, msg, ph)
				}
			case "between":
				if msg, ph, hasError := Between(value, params...); hasError {
					e.Append(tagName, msg, ph)
				}
			case "regex":
				if msg, ph, hasError := Regex(value, params...); hasError {
					e.Append(tagName, msg, ph)
				}
			case "exists":
				if msg, ph, hasError := Exists(value, params...); hasError {
					e.Append(tagName, msg, ph)
				}
			case "not_exists", "notExists":
				if msg, ph, hasError := NotExists(value, params...); hasError {
					e.Append(tagName, msg, ph)
				}
			case "unique":
				params = append(params, c.Request.URL.Query().Get("id"))
				if msg, ph, hasError := Unique(value, params...); hasError {
					e.Append(tagName, msg, ph)
				}
			case "in":
				if msg, ph, hasError := In(value, params...); hasError {
					e.Append(tagName, msg, ph)
				}
			case "nin":
				if msg, ph, hasError := Nin(value, params...); hasError {
					e.Append(tagName, msg, ph)
				}
			case "before":
				if msg, ph, hasError := Before(value, params...); hasError {
					e.Append(tagName, msg, ph)
				}
			case "after":
				if msg, ph, hasError := After(value, params...); hasError {
					e.Append(tagName, msg, ph)
				}
			case "before_now", "beforeNow":
				if msg, ph, hasError := BeforeNow(value, params...); hasError {
					e.Append(tagName, msg, ph)
				}
			case "after_now", "afterNow":
				if msg, ph, hasError := AfterNow(value, params...); hasError {
					e.Append(tagName, msg, ph)
				}
			}
		}
	}
	return e.HasErrors()
}
