package err

import (
	"donbarrigon/new/internal/utils/fm"
	"donbarrigon/new/internal/utils/lang"
	"encoding/json"
)

type ValidationError struct {
	Messages     map[string][]string
	Placeholders map[string][]fm.Placeholder
	IsEntity     bool
}

func NewValidationError() *ValidationError {
	return &ValidationError{
		Messages:     map[string][]string{},
		Placeholders: map[string][]fm.Placeholder{},
	}
}

func (e *ValidationError) Append(field string, message string, ph fm.Placeholder) {
	ph.Append("field", field)
	e.Messages[field] = append(e.Messages[field], message)
	e.Placeholders[field] = append(e.Placeholders[field], ph)
}

func (e *ValidationError) AppendM(field string, message string) {
	e.Messages[field] = append(e.Messages[field], message)
	e.Placeholders[field] = append(e.Placeholders[field], fm.Placeholder{"field": field})
}

func (e *ValidationError) HasErrors() error {
	if len(e.Messages) > 0 {
		return e
	}
	return nil
}

// ================================
// Funciones para la interfaz de Errror
// ================================

func (e *ValidationError) Error() string {
	b, err := json.Marshal(e.Messages)
	if err != nil {
		return "{}"
	}
	return string(b)
}
func (e *ValidationError) Herror(l string) *HttpError {
	result := map[string][]string{}
	if len(e.Messages) > 0 {
		for key, messages := range e.Messages {
			result[key] = []string{}
			for i, message := range messages {
				result[key] = append(result[key], lang.T(l, message, e.Placeholders[key][i]))
			}
		}
	}
	if e.IsEntity {
		return &HttpError{Status: UNPROCESSABLE_ENTITY, Message: "No pudimos procesar la información que enviaste", Err: result}
	}
	return &HttpError{Status: BAD_REQUEST, Message: "Algo no está bien con tu solicitud", Err: result}
}
