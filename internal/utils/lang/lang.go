package lang

import (
	"donbarrigon/new/internal/utils/str"
)

func T(lang string, text string, ph str.Placeholder) string {
	txt := MessagesMap[lang][text]
	if txt == "" {
		txt = text
	}

	if ph == nil {
		return txt
	}

	for _, p := range ph {
		v := AtributesMap[lang][p.Value]
		if v != "" {
			ph.Rename(p.Key, v)
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
