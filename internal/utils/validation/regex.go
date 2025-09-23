package validation

var RegexPatterns = map[string]string{
	"alpha":                         `^[A-Za-z]+$`,
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
	"phone":                         `^\+?\d{1,4}?[-.\s]?\(?\d{1,4}\)?([-.\s]?\d{1,9}){1,3}$`,
	"username":                      `^[A-Za-zÑñ][^\p{Z}\p{M}\p{So}\p{Sk}\p{Lm}]*$`,
	"emial":                         `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[A-Za-z]{2,}$`,
}
