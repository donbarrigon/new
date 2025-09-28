package validate

import (
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/fm"
	"donbarrigon/new/internal/utils/handler"
	"reflect"
	"strings"
)

type ValidationFunc func(value reflect.Value, params ...string) (string, fm.Placeholder, bool)

type Rule struct {
	Fun    ValidationFunc
	Params []string
}

type Rules map[string][]Rule

type Validator interface {
	Rules() Rules
	PrepareForValidation(c *handler.HttpContext) *err.ValidationError
}

func Body(c *handler.HttpContext, validator Validator) err.Error {
	if e := c.GetBody(validator); e != nil {
		return e
	}
	var e *err.ValidationError
	e = validator.PrepareForValidation(c)
	if e == nil {
		e = err.NewValidationError()
	}
	rules := validator.Rules()
	return validate(validator, rules, e)
}

func From(validator Validator) err.Error {
	var e *err.ValidationError
	e = validator.PrepareForValidation(nil)
	if e == nil {
		e = err.NewValidationError()
	}
	return validate(validator, validator.Rules(), e)
}

func Struct(data any, rules Rules) err.Error {
	var e *err.ValidationError
	e = err.NewValidationError()
	return validate(data, rules, e)
}

func validate(validator any, rules Rules, e *err.ValidationError) err.Error {
	val := reflect.ValueOf(validator)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return err.Internal("Error de validación: Los datos tienen que ser un *struct")
	}
	typ := val.Type()
	numFields := typ.NumField()
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		tagName := strings.Split(tag, ",")[0]
		if tagName == "-" || tagName == "id" {
			continue
		}
		if tagName == "" {
			tagName = field.Name
		}
		validations := rules[tagName]
		if validations == nil {
			continue
		}
		value := val.Field(i)

		for _, rule := range validations {
			if msg, ph, hasError := rule.Fun(value, rule.Params...); hasError {
				e.Append(tagName, msg, ph)
			}
		}
	}
	return e.HasErrors()
}
