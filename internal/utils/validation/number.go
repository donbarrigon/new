package validation

import (
	"donbarrigon/new/internal/utils/fm"
	"reflect"
	"strconv"
	"time"
)

// Assuming fm.fm.Placeholder is defined somewhere like this:
// type fm.Placeholder map[string]string

// Min valida que el valor sea mayor o igual al mínimo especificado
func Min(value reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	if len(params) < 1 {
		return "Parámetro mínimo requerido", fm.Placeholder{}, true
	}

	minStr := params[0]
	ph := fm.Placeholder{"min": minStr}

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

	return "", nil, false
}

// Max valida que el valor sea menor o igual al máximo especificado
func Max(value reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	if len(params) < 1 {
		return "Parámetro máximo requerido", fm.Placeholder{}, true
	}

	maxStr := params[0]
	ph := fm.Placeholder{"max": maxStr}

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

	return "", nil, false
}

// Between valida que el valor esté dentro del rango especificado
func Between(value reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	if len(params) < 2 {
		return "Se requieren parámetros mínimo y máximo", fm.Placeholder{}, true
	}

	minStr := params[0]
	maxStr := params[1]
	ph := fm.Placeholder{"min": minStr, "max": maxStr}

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

	return "", nil, false
}
