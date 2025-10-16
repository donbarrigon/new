package validate

import (
	"donbarrigon/new/internal/utils/fm"
	"reflect"
	"time"
)

// Before valida que la fecha sea anterior a la fecha especificada
func Before(value reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	if len(params) < 1 {
		return "Parámetro de fecha requerido", fm.Placeholder{}, true
	}

	beforeStr := params[0]
	ph := fm.Placeholder{"date": beforeStr}

	if value.Type() != reflect.TypeOf(time.Time{}) {
		return "El campo :field debe ser de tipo fecha", ph, true
	}

	before, err := time.Parse("2006-01-02", beforeStr)
	if err != nil {
		before, err = time.Parse("2006-01-02T15:04:05Z07:00", beforeStr)
		if err != nil {
			return "Formato de fecha inválido", ph, true
		}
	}

	fieldTime := value.Interface().(time.Time)
	if !fieldTime.Before(before) {
		return "El campo :field debe ser anterior a :date", ph, true
	}

	return "", nil, false
}

// After valida que la fecha sea posterior a la fecha especificada
func After(value reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	if len(params) < 1 {
		return "Parámetro de fecha requerido", fm.Placeholder{}, true
	}

	afterStr := params[0]
	ph := fm.Placeholder{"date": afterStr}

	if value.Type() != reflect.TypeOf(time.Time{}) {
		return "El campo :field debe ser de tipo fecha", ph, true
	}

	after, err := time.Parse("2006-01-02", afterStr)
	if err != nil {
		after, err = time.Parse("2006-01-02T15:04:05Z07:00", afterStr)
		if err != nil {
			return "Formato de fecha inválido", ph, true
		}
	}

	fieldTime := value.Interface().(time.Time)
	if !fieldTime.After(after) {
		return "El campo :field debe ser posterior a :date", ph, true
	}

	return "", nil, false
}

// BeforeNow valida que la fecha sea anterior a la fecha actual
func BeforeNow(value reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	ph := fm.Placeholder{"date": "ahora"}

	if value.Type() != reflect.TypeOf(time.Time{}) {
		return "El campo :field debe ser de tipo fecha", ph, true
	}

	fieldTime := value.Interface().(time.Time)
	now := time.Now()

	if !fieldTime.Before(now) {
		return "El campo :field debe ser anterior a la fecha actual", ph, true
	}

	return "", nil, false
}

// AfterNow valida que la fecha sea posterior a la fecha actual
func AfterNow(value reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	ph := fm.Placeholder{"date": "ahora"}

	if value.Type() != reflect.TypeOf(time.Time{}) {
		return "El campo :field debe ser de tipo fecha", ph, true
	}

	fieldTime := value.Interface().(time.Time)
	now := time.Now()

	if !fieldTime.After(now) {
		return "El campo :field debe ser posterior a la fecha actual", ph, true
	}

	return "", nil, false
}
