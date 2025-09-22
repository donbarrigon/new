package err

import (
	"donbarrigon/new/internal/utils/fm"
	"donbarrigon/new/internal/utils/lang"
)

type ValidationError struct {
	Messages     map[string][]string
	Placeholders map[string][]fm.Placeholder
}

func NewValidationError() *ValidationError {
	return &ValidationError{
		Messages:     map[string][]string{},
		Placeholders: map[string][]fm.Placeholder{},
	}
}

func (e *ValidationError) Append(field string, message string, ph fm.Placeholder) {
	e.Messages[field] = append(e.Messages[field], message)
	e.Placeholders[field] = append(e.Placeholders[field], ph)
}

func (e *ValidationError) AppendM(field string, message string) {
	e.Messages[field] = append(e.Messages[field], message)
	e.Placeholders[field] = append(e.Placeholders[field], fm.Placeholder{})
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Messages) > 0
}

func (e *ValidationError) Errors(l string) map[string][]string {
	result := map[string][]string{}

	for key, messages := range e.Messages {
		result[key] = []string{}
		for i, message := range messages {
			result[key] = append(result[key], lang.T(l, message, e.Placeholders[key][i]))
		}
	}
	return result
}
