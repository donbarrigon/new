package validation

import (
	"donbarrigon/new/internal/utils/str"
	"reflect"
	"regexp"
)

var RegexMap = map[string]string{
	"snake_case":                    "^[a-z][a-z0-9_]*$",
	"kebab-case":                    "^[a-z][a-z0-9-]*$",
	"camelCase":                     "^[a-z]+(?:[A-Z][a-z0-9]*)*$",
	"PascalCase":                    "^[A-Z][a-z0-9]*(?:[A-Z][a-z0-9]*)*$",
	"alpha":                         `^[A-Za-z]+$`,
	"numeric":                       `^[0-9]+$`,
	"alpha_dash":                    `^[A-Za-z-_]+$`,
	"alpha_spaces":                  `^[A-Za-z ]+$`,
	"alpha_dash_spaces":             `^[A-Za-z -_]+$`,
	"alpha_num":                     `^[A-Za-z0-9]+$`,
	"alpha_num_dash":                `^[A-Za-z0-9-_]+$`,
	"alpha_num_spaces":              `^[A-Za-z0-9 ]+$`,
	"alpha_num_dash_spaces":         `^[A-Za-z0-9 -_]+$`,
	"alpha_accents":                 `^[A-Za-zÁÉÍÓÚÜÑáéíóúüñ]+$`,
	"alpha_dash_accents":            `^[A-Za-zÁÉÍÓÚÜÑáéíóúüñ-_]+$`,
	"alpha_spaces_accents":          `^[A-Za-zÁÉÍÓÚÜÑáéíóúüñ ]+$`,
	"alpha_dash_spaces_accents":     `^[A-Za-zÁÉÍÓÚÜÑáéíóúüñ -_]+$`,
	"alpha_num_accents":             `^[A-Za-zÁÉÍÓÚÜÑáéíóúüñ0-9]+$`,
	"alpha_num_dash_accents":        `^[A-Za-zÁÉÍÓÚÜÑáéíóúüñ0-9-_]+$`,
	"alpha_num_spaces_accents":      `^[A-Za-zÁÉÍÓÚÜÑáéíóúüñ0-9 ]+$`,
	"alpha_num_dash_spaces_accents": `^[A-Za-zÁÉÍÓÚÜÑáéíóúüñ0-9 -_]+$`,
	"nickname":                      `^[A-Za-zÑñ][^\p{Z}\p{M}\p{So}\p{Sk}\p{Lm}]*$`,
	"username":                      `^[a-zA-Z0-9_]{3,20}$`,
	"email":                         `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
	"phone":                         `^[\+]?[1-9][\d]{0,15}$`,
	"url":                           `^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$`,
	"uuid":                          `^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
	"ipv4":                          `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`,
	"mac":                           `^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`,
	"creditcard":                    `^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|3[47][0-9]{13}|3[0-9]{13}|6(?:011|5[0-9]{2})[0-9]{12})$`,
	"date":                          `^\d{4}-\d{2}-\d{2}$`,
	"time":                          `^([01]?[0-9]|2[0-3]):[0-5][0-9]$`,
	"datetime":                      `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{3})?Z?$`,
	"password":                      `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`,
	"slug":                          `^[a-z0-9]+(?:-[a-z0-9]+)*$`,
	"hexcolor":                      `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`,
	"base64":                        `^[A-Za-z0-9+/]*={0,2}$`,
	"jwt":                           `^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$`,
	"isbn":                          `^(?:ISBN(?:-1[03])?:? )?(?=[0-9X]{10}$|(?=(?:[0-9]+[- ]){3})[- 0-9X]{13}$|97[89][0-9]{10}$|(?=(?:[0-9]+[- ]){4})[- 0-9]{17}$)(?:97[89][- ]?)?[0-9]{1,5}[- ]?[0-9]+[- ]?[0-9]+[- ]?[0-9X]$`,
}

var regexCache = map[string]*regexp.Regexp{}

// Regex valida que el valor coincida con la expresión regular especificada
func Regex(value reflect.Value, params ...string) (string, str.Placeholder, bool) {
	if len(params) < 1 {
		return "Se requiere un patrón de expresión regular", str.Placeholder{}, true
	}

	// Solo validamos strings
	if value.Kind() != reflect.String {
		return "El campo :field debe ser de tipo string para validación regex", str.Placeholder{}, true
	}

	valueStr := value.String()
	regexPattern := params[0]
	ph := str.Placeholder{{Key: "regex", Value: regexPattern}}

	// Buscar en el mapa de regex predefinidas
	compiledRegex := regexCache[regexPattern]
	if compiledRegex == nil {
		var e error
		compiledRegex, e = regexp.Compile(regexPattern)
		if e != nil {
			return "Patrón de expresión regular inválido [:regex]", ph, true
		}
		regexCache[regexPattern] = compiledRegex
	}

	if !compiledRegex.MatchString(valueStr) {
		if len(params) > 1 && params[1] != "" {
			return params[1], ph, true
		}
		return "El campo :field no tiene el formato correcto", ph, true
	}

	return "", nil, false
}
