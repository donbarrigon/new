package validation

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Assuming fm.Placeholder is defined somewhere like this:
// type Placeholder map[string]string

type Placeholder map[string]string

// Min valida que el valor sea mayor o igual al mínimo especificado
func Min(value reflect.Value, params ...string) (string, Placeholder, bool) {
	if len(params) < 1 {
		return "Parámetro mínimo requerido", Placeholder{}, true
	}

	minStr := params[0]
	ph := Placeholder{"min": minStr}

	switch value.Kind() {
	case reflect.String:
		minLen, err := strconv.Atoi(minStr)
		if err != nil {
			return "Parámetro mínimo inválido", ph, true
		}
		if len(value.String()) < minLen {
			return "El campo :field requiere mínimo :min caracteres", ph, true
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil {
			return "Parámetro mínimo inválido", ph, true
		}
		if value.Int() < min {
			return "El campo :field debe ser mayor o igual a :min", ph, true
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		min, err := strconv.ParseUint(minStr, 10, 64)
		if err != nil {
			return "Parámetro mínimo inválido", ph, true
		}
		if value.Uint() < min {
			return "El campo :field debe ser mayor o igual a :min", ph, true
		}

	case reflect.Float32, reflect.Float64:
		min, err := strconv.ParseFloat(minStr, 64)
		if err != nil {
			return "Parámetro mínimo inválido", ph, true
		}
		if value.Float() < min {
			return "El campo :field debe ser mayor o igual a :min", ph, true
		}

	case reflect.Slice, reflect.Array:
		minLen, err := strconv.Atoi(minStr)
		if err != nil {
			return "Parámetro mínimo inválido", ph, true
		}
		if value.Len() < minLen {
			return "El campo :field debe tener mínimo :min elementos", ph, true
		}

	default:
		// Para fechas (time.Time)
		if value.Type() == reflect.TypeOf(time.Time{}) {
			min, err := time.Parse("2006-01-02", minStr)
			if err != nil {
				min, err = time.Parse("2006-01-02T15:04:05Z07:00", minStr)
				if err != nil {
					return "Formato de fecha mínima inválido", ph, true
				}
			}
			fieldTime := value.Interface().(time.Time)
			if fieldTime.Before(min) {
				return "El campo :field debe ser posterior o igual a :min", ph, true
			}
		} else {
			return "Tipo no soportado para validación Min", ph, true
		}
	}

	return "", ph, false
}

// Max valida que el valor sea menor o igual al máximo especificado
func Max(value reflect.Value, params ...string) (string, Placeholder, bool) {
	if len(params) < 1 {
		return "Parámetro máximo requerido", Placeholder{}, true
	}

	maxStr := params[0]
	ph := Placeholder{"max": maxStr}

	switch value.Kind() {
	case reflect.String:
		maxLen, err := strconv.Atoi(maxStr)
		if err != nil {
			return "Parámetro máximo inválido", ph, true
		}
		if len(value.String()) > maxLen {
			return "El campo :field no puede exceder :max caracteres", ph, true
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		max, err := strconv.ParseInt(maxStr, 10, 64)
		if err != nil {
			return "Parámetro máximo inválido", ph, true
		}
		if value.Int() > max {
			return "El campo :field debe ser menor o igual a :max", ph, true
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		max, err := strconv.ParseUint(maxStr, 10, 64)
		if err != nil {
			return "Parámetro máximo inválido", ph, true
		}
		if value.Uint() > max {
			return "El campo :field debe ser menor o igual a :max", ph, true
		}

	case reflect.Float32, reflect.Float64:
		max, err := strconv.ParseFloat(maxStr, 64)
		if err != nil {
			return "Parámetro máximo inválido", ph, true
		}
		if value.Float() > max {
			return "El campo :field debe ser menor o igual a :max", ph, true
		}

	case reflect.Slice, reflect.Array:
		maxLen, err := strconv.Atoi(maxStr)
		if err != nil {
			return "Parámetro máximo inválido", ph, true
		}
		if value.Len() > maxLen {
			return "El campo :field no puede tener más de :max elementos", ph, true
		}

	default:
		// Para fechas (time.Time)
		if value.Type() == reflect.TypeOf(time.Time{}) {
			max, err := time.Parse("2006-01-02", maxStr)
			if err != nil {
				max, err = time.Parse("2006-01-02T15:04:05Z07:00", maxStr)
				if err != nil {
					return "Formato de fecha máxima inválido", ph, true
				}
			}
			fieldTime := value.Interface().(time.Time)
			if fieldTime.After(max) {
				return "El campo :field debe ser anterior o igual a :max", ph, true
			}
		} else {
			return "Tipo no soportado para validación Max", ph, true
		}
	}

	return "", ph, false
}

// Between valida que el valor esté dentro del rango especificado
func Between(value reflect.Value, params ...string) (string, Placeholder, bool) {
	if len(params) < 2 {
		return "Se requieren parámetros mínimo y máximo", Placeholder{}, true
	}

	minStr := params[0]
	maxStr := params[1]
	ph := Placeholder{"min": minStr, "max": maxStr}

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min, err1 := strconv.ParseInt(minStr, 10, 64)
		max, err2 := strconv.ParseInt(maxStr, 10, 64)
		if err1 != nil || err2 != nil {
			return "Parámetros de rango inválidos", ph, true
		}
		val := value.Int()
		if val < min || val > max {
			return "El campo :field debe estar entre :min y :max", ph, true
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		min, err1 := strconv.ParseUint(minStr, 10, 64)
		max, err2 := strconv.ParseUint(maxStr, 10, 64)
		if err1 != nil || err2 != nil {
			return "Parámetros de rango inválidos", ph, true
		}
		val := value.Uint()
		if val < min || val > max {
			return "El campo :field debe estar entre :min y :max", ph, true
		}

	case reflect.Float32, reflect.Float64:
		min, err1 := strconv.ParseFloat(minStr, 64)
		max, err2 := strconv.ParseFloat(maxStr, 64)
		if err1 != nil || err2 != nil {
			return "Parámetros de rango inválidos", ph, true
		}
		val := value.Float()
		if val < min || val > max {
			return "El campo :field debe estar entre :min y :max", ph, true
		}

	default:
		// Para fechas (time.Time)
		if value.Type() == reflect.TypeOf(time.Time{}) {
			min, err1 := time.Parse("2006-01-02", minStr)
			if err1 != nil {
				min, err1 = time.Parse("2006-01-02T15:04:05Z07:00", minStr)
			}
			max, err2 := time.Parse("2006-01-02", maxStr)
			if err2 != nil {
				max, err2 = time.Parse("2006-01-02T15:04:05Z07:00", maxStr)
			}
			if err1 != nil || err2 != nil {
				return "Formato de fechas de rango inválido", ph, true
			}
			fieldTime := value.Interface().(time.Time)
			if fieldTime.Before(min) || fieldTime.After(max) {
				return "El campo :field debe estar entre :min y :max", ph, true
			}
		} else {
			return "Tipo no soportado para validación Between", ph, true
		}
	}

	return "", ph, false
}

// Before valida que la fecha sea anterior a la fecha especificada
func Before(value reflect.Value, params ...string) (string, Placeholder, bool) {
	if len(params) < 1 {
		return "Parámetro de fecha requerido", Placeholder{}, true
	}

	beforeStr := params[0]
	ph := Placeholder{"date": beforeStr}

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

	return "", ph, false
}

// After valida que la fecha sea posterior a la fecha especificada
func After(value reflect.Value, params ...string) (string, Placeholder, bool) {
	if len(params) < 1 {
		return "Parámetro de fecha requerido", Placeholder{}, true
	}

	afterStr := params[0]
	ph := Placeholder{"date": afterStr}

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

	return "", ph, false
}

// BeforeNow valida que la fecha sea anterior a la fecha actual
func BeforeNow(value reflect.Value, params ...string) (string, Placeholder, bool) {
	ph := Placeholder{"date": "ahora"}

	if value.Type() != reflect.TypeOf(time.Time{}) {
		return "El campo :field debe ser de tipo fecha", ph, true
	}

	fieldTime := value.Interface().(time.Time)
	now := time.Now()

	if !fieldTime.Before(now) {
		return "El campo :field debe ser anterior a la fecha actual", ph, true
	}

	return "", ph, false
}

// AfterNow valida que la fecha sea posterior a la fecha actual
func AfterNow(value reflect.Value, params ...string) (string, Placeholder, bool) {
	ph := Placeholder{"date": "ahora"}

	if value.Type() != reflect.TypeOf(time.Time{}) {
		return "El campo :field debe ser de tipo fecha", ph, true
	}

	fieldTime := value.Interface().(time.Time)
	now := time.Now()

	if !fieldTime.After(now) {
		return "El campo :field debe ser posterior a la fecha actual", ph, true
	}

	return "", ph, false
}

// In valida que el valor esté presente en la lista de valores permitidos
func In(value reflect.Value, params ...string) (string, Placeholder, bool) {
	if len(params) < 1 {
		return "Se requiere al menos un valor permitido", Placeholder{}, true
	}

	ph := Placeholder{"values": fmt.Sprintf("%v", params)}

	switch value.Kind() {
	case reflect.String:
		valueStr := value.String()
		for _, param := range params {
			if valueStr == param {
				return "", ph, false
			}
		}
		return "El campo :field debe ser uno de los valores :values", ph, true

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valueInt := value.Int()
		for _, param := range params {
			if paramInt, e := strconv.ParseInt(param, 10, 64); e == nil {
				if valueInt == paramInt {
					return "", ph, false
				}
			}
		}
		return "El campo :field debe ser uno de los valores :values", ph, true

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valueUint := value.Uint()
		for _, param := range params {
			if paramUint, e := strconv.ParseUint(param, 10, 64); e == nil {
				if valueUint == paramUint {
					return "", ph, false
				}
			}
		}
		return "El campo :field debe ser uno de los valores :values", ph, true

	case reflect.Float32, reflect.Float64:
		valueFloat := value.Float()
		for _, param := range params {
			if paramFloat, e := strconv.ParseFloat(param, 64); e == nil {
				if valueFloat == paramFloat {
					return "", ph, false
				}
			}
		}
		return "El campo :field debe ser uno de los valores :values", ph, true

	default:
		return "Tipo no soportado para validación In", ph, true
	}
}

// Nin valida que el valor NO esté presente en la lista de valores prohibidos
func Nin(value reflect.Value, params ...string) (string, Placeholder, bool) {
	if len(params) < 1 {
		return "Se requiere al menos un valor prohibido", Placeholder{}, true
	}

	ph := Placeholder{"values": fmt.Sprintf("%v", params)}

	switch value.Kind() {
	case reflect.String:
		valueStr := value.String()
		for _, param := range params {
			if valueStr == param {
				return "El campo :field no puede ser uno de los valores :values", ph, true
			}
		}
		return "", ph, false

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valueInt := value.Int()
		for _, param := range params {
			if paramInt, e := strconv.ParseInt(param, 10, 64); e == nil {
				if valueInt == paramInt {
					return "El campo :field no puede ser uno de los valores :values", ph, true
				}
			}
		}
		return "", ph, false

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valueUint := value.Uint()
		for _, param := range params {
			if paramUint, e := strconv.ParseUint(param, 10, 64); e == nil {
				if valueUint == paramUint {
					return "El campo :field no puede ser uno de los valores :values", ph, true
				}
			}
		}
		return "", ph, false

	case reflect.Float32, reflect.Float64:
		valueFloat := value.Float()
		for _, param := range params {
			if paramFloat, e := strconv.ParseFloat(param, 64); e == nil {
				if valueFloat == paramFloat {
					return "El campo :field no puede ser uno de los valores :values", ph, true
				}
			}
		}
		return "", ph, false

	default:
		return "Tipo no soportado para validación Nin", ph, true
	}
}
