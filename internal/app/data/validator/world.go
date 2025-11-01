package validator

import (
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/handler"
	"donbarrigon/new/internal/utils/str"
	"donbarrigon/new/internal/utils/validation"
	"regexp"
)

// ================================================================
//                         COUNTRIES
// ================================================================

type CountryStore struct {
	Name           string                 `json:"name,omitempty"`
	Iso3           string                 `json:"iso3,omitempty"`
	Iso2           string                 `json:"iso2,omitempty"`
	NumericCode    string                 `json:"numeric_code,omitempty"`
	PhoneCode      string                 `json:"phonecode,omitempty"`
	Capital        string                 `json:"capital,omitempty"`
	Currency       string                 `json:"currency,omitempty"`
	CurrencyName   string                 `json:"currencyName,omitempty"`
	CurrencySymbol string                 `json:"currencySymbol,omitempty"`
	TLD            string                 `json:"tld,omitempty"`
	Native         string                 `json:"native,omitempty"`
	Region         CountryStoreRegion     `json:"region"`
	Subregion      CountryStoreSubRegion  `json:"subregion"`
	Nationality    string                 `json:"nationality,omitempty"`
	Timezones      []CountryStoreTimezone `json:"timezones,omitempty"`
	Translations   map[string]string      `json:"translations,omitempty"`
	Location       db.GeoPoint            `json:"location"`
	Emoji          string                 `json:"emoji,omitempty"`
	EmojiU         string                 `json:"emojiU,omitempty"`
}

type CountryStoreTimezone struct {
	ZoneName      string `json:"zoneName,omitempty"`
	GMTOffset     int    `json:"gmtOffset,omitempty"`
	GMTOffsetName string `json:"gmtOffsetName,omitempty"`
	Abbreviation  string `json:"abbreviation,omitempty"`
	TZName        string `json:"tzName,omitempty"`
}

type CountryStoreRegion struct {
	ID           int               `json:"id,omitempty"`
	Name         string            `json:"name,omitempty"`
	Translations map[string]string `json:"translations,omitempty"`
	WikiDataId   string            `json:"wikiDataId,omitempty"`
}

type CountryStoreSubRegion struct {
	ID           int               `json:"id,omitempty"`
	RegionID     int               `json:"regionId,omitempty"`
	Name         string            `json:"name,omitempty"`
	Translations map[string]string `json:"translations,omitempty"`
	WikiDataId   string            `json:"wikiDataId,omitempty"`
}

func (cs *CountryStore) Rules() validation.Rules {
	return validation.Rules{
		"name": {
			"required": {},
			"between":  {"2", "100"},
			"unique":   {"countries", "name"},
		},
		"iso3": {
			"required": {},
			"regex":    {"^[A-Z]{3}$"},
			"unique":   {"countries", "iso3"},
		},
		"iso2": {
			"required": {},
			"regex":    {"^[A-Z]{2}$"},
			"unique":   {"countries", "iso2"},
		},
		"numeric_code": {
			"required": {},
			"regex":    {"^[0-9]{3}$"},
			"unique":   {"countries", "numeric_code"},
		},
		"phonecode": {
			"required": {},
			"regex":    {"^[0-9]+$"},
			"between":  {"1", "6"},
		},
		"capital": {
			"required": {},
			"between":  {"2", "100"},
		},
		"currency": {
			"required": {},
			"regex":    {"^[A-Z]{3}$"},
		},
		"currencyName": {
			"required": {},
			"between":  {"2", "100"},
		},
		"currencySymbol": {
			"between": {"1", "5"},
		},
		"tld": {
			"regex":   {"^\\.[a-z]{2,}$"},
			"between": {"2", "10"},
		},
		"native": {
			"between": {"2", "100"},
		},
		"nationality": {
			"between": {"2", "100"},
		},
		"timezones": {
			"min_items": {"1"},
		},
		"timezones.*.zoneName": {
			"required": {},
			"between":  {"2", "100"},
		},
		"timezones.*.gmtOffset": {
			"required": {},
			"numeric":  {},
		},
		"timezones.*.abbreviation": {
			"between": {"1", "10"},
		},
		"timezones.*.tzName": {
			"between": {"2", "100"},
		},
		"translations": {
			"map": {"string", "string"},
		},
		"location": {
			"required":  {},
			"geo_point": {},
		},
		"emoji": {
			"regex": {"^[\u263a-\U0001f645]+$"},
		},
		"emojiU": {
			"regex": {"^U\\+[A-F0-9]{4,6}$"},
		},
	}
}

func (cs *CountryStore) PrepareForValidation(c *handler.Context) *err.ValidationError {
	e := err.NewValidationError()
	// ----------------------------------------------------------------
	// Validaciones de la region
	// ----------------------------------------------------------------
	if cs.Region.ID == 0 {
		e.AppendM("region.id", "El ID de la región es requerido")
	}
	if len(cs.Region.Name) < 2 || len(cs.Region.Name) > 255 {
		e.AppendM("region.name", "El nombre de la región debe tener entre 2 y 255 caracteres")
	}
	srx := "^[A-Za-z0-9]+$"
	crx, er := regexp.Compile(srx)
	if er != nil {
		e.Append("region.wikiDataId", "Patrón de expresión regular inválido [:field / :regex]", str.NewPlaceholder("regex", srx))
	}
	if !crx.MatchString(cs.Region.WikiDataId) {
		e.AppendM("region.wikiDataId", "El wikiDataId de la región es inválido")
	}
	if cs.Region.Translations == nil {
		cs.Region.Translations = map[string]string{}
	}

	// ----------------------------------------------------------------
	// Validaciones de la subregion
	// ----------------------------------------------------------------
	if cs.Subregion.ID == 0 {
		e.AppendM("subregion.id", "El ID de la subregión es requerido")
	}
	if cs.Subregion.RegionID == 0 {
		e.AppendM("subregion.regionId", "El ID de la región es requerido")
	}
	if len(cs.Subregion.Name) < 2 || len(cs.Subregion.Name) > 255 {
		e.AppendM("subregion.name", "El nombre de la subregión debe tener entre 2 y 255 caracteres")
	}
	if !crx.MatchString(cs.Subregion.WikiDataId) {
		e.AppendM("subregion.wikiDataId", "El wikiDataId de la subregión es inválido")
	}
	if cs.Subregion.Translations == nil {
		cs.Subregion.Translations = map[string]string{}
	}

	// ----------------------------------------------------------------
	// Validaciones de la subregion
	// ----------------------------------------------------------------
	if cs.Timezones == nil {
		cs.Timezones = []CountryStoreTimezone{}
	}
	return e
}

// ================================================================
//                         STATES
// ================================================================
