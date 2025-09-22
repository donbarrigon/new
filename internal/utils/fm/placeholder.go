package fm

import (
	"fmt"
	"strings"
)

type Placeholder map[string]string

func (ph Placeholder) Append(key string, value any) {
	ph[key] = fmt.Sprintf("%v", value)
}

func (ph Placeholder) Remove(key string) {
	delete(ph, key)
}

// InteporlatePlaceholder Reemplaza los placeholders en el texto
func (ph Placeholder) Replace(text string) string {
	for k, v := range ph {
		text = strings.ReplaceAll(text, ":"+k, v)
	}
	return text
}
