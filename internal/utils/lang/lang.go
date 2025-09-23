package lang

import (
	"donbarrigon/new/internal/utils/fm"
)

func T(lang string, text string, ph fm.Placeholder) string {
	txt := MessagesMap[lang][text]
	if txt == "" {
		txt = text
	}

	if ph == nil {
		return txt
	}

	for key, value := range ph {
		v := AtributesMap[lang][value]
		if v != "" {
			ph[key] = v
		}
	}
	return ph.Replace(txt)
}

var MessagesMap = map[string]map[string]string{
	"es": {},
}

var AtributesMap = map[string]map[string]string{
	"es": {},
}
