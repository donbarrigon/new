package model

import (
	"context"
	"donbarrigon/new/internal/app/data/validator"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/err"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ================================================================
// countries
// ================================================================
type Country struct {
	ID             bson.ObjectID     `bson:"id,omitempty"              json:"id,omitempty"`
	Name           string            `bson:"name,omitempty"            json:"name,omitempty"`
	Iso3           string            `bson:"iso3,omitempty"            json:"iso3,omitempty"`
	Iso2           string            `bson:"iso2,omitempty"            json:"iso2,omitempty"`
	NumericCode    string            `bson:"numeric_code,omitempty"    json:"numeric_code,omitempty"`
	PhoneCode      string            `bson:"phonecode,omitempty"       json:"phonecode,omitempty"`
	Capital        string            `bson:"capital,omitempty"         json:"capital,omitempty"`
	Currency       string            `bson:"currency,omitempty"        json:"currency,omitempty"`
	CurrencyName   string            `bson:"currency_name,omitempty"   json:"currencyName,omitempty"`
	CurrencySymbol string            `bson:"currency_symbol,omitempty" json:"currencySymbol,omitempty"`
	TLD            string            `bson:"tld,omitempty"             json:"tld,omitempty"`
	Native         string            `bson:"native,omitempty"          json:"native,omitempty"`
	Region         CountryRegion     `bson:"region"                    json:"region"`
	Subregion      CountrySubRegion  `bson:"subregion"                 json:"subregion"`
	Nationality    string            `bson:"nationality,omitempty"     json:"nationality,omitempty"`
	Timezones      []CountryTimezone `bson:"timezones,omitempty"       json:"timezones,omitempty"`
	Translations   map[string]string `bson:"translations,omitempty"    json:"translations,omitempty"`
	Location       db.GeoPoint       `bson:"location"                  json:"location"`
	Emoji          string            `bson:"emoji,omitempty"           json:"emoji,omitempty"`
	EmojiU         string            `bson:"emojiU,omitempty"          json:"emojiU,omitempty"`
	CreatedAt      time.Time         `bson:"created_at"                json:"createdAt"`
	UpdatedAt      time.Time         `bson:"updated_at"                json:"updatedAt"`
}

type CountryTimezone struct {
	ZoneName      string `bson:"zone_name,omitempty"       json:"zoneName,omitempty"`
	GMTOffset     int    `bson:"gmt_offset,omitempty"      json:"gmtOffset,omitempty"`
	GMTOffsetName string `bson:"gmt_offset_name,omitempty" json:"gmtOffsetName,omitempty"`
	Abbreviation  string `bson:"abbreviation,omitempty"    json:"abbreviation,omitempty"`
	TZName        string `bson:"tz_name,omitempty"         json:"tzName,omitempty"`
}

type CountryRegion struct {
	ID           int               `bson:"id,omitempty"           json:"id,omitempty"`
	Name         string            `bson:"name,omitempty"         json:"name,omitempty"`
	Translations map[string]string `bson:"translations,omitempty" json:"translations,omitempty"`
	WikiDataId   string            `bson:"wiki_data_id,omitempty" json:"wikiDataId,omitempty"`
}

type CountrySubRegion struct {
	ID           int               `bson:"id,omitempty"           json:"id,omitempty"`
	RegionID     int               `bson:"region_id,omitempty"    json:"regionId,omitempty"`
	Name         string            `bson:"name,omitempty"         json:"name,omitempty"`
	Translations map[string]string `bson:"translations,omitempty" json:"translations,omitempty"`
	WikiDataId   string            `bson:"wiki_data_id,omitempty" json:"wikiDataId,omitempty"`
}

func (c *Country) Coll() string         { return "countries" }
func (c *Country) GetID() bson.ObjectID { return c.ID }

func CountryCreate(dto *validator.CountryStore) (*Country, error) {

	region := CountryRegion{
		ID:           dto.Region.ID,
		Name:         dto.Region.Name,
		Translations: dto.Region.Translations,
		WikiDataId:   dto.Region.WikiDataId,
	}
	Subregion := CountrySubRegion{
		ID:           dto.Subregion.ID,
		RegionID:     dto.Subregion.RegionID,
		Name:         dto.Subregion.Name,
		Translations: dto.Subregion.Translations,
		WikiDataId:   dto.Subregion.WikiDataId,
	}
	timezones := []CountryTimezone{}
	for _, v := range dto.Timezones {
		ct := CountryTimezone{
			ZoneName:      v.ZoneName,
			GMTOffset:     v.GMTOffset,
			GMTOffsetName: v.GMTOffsetName,
			Abbreviation:  v.Abbreviation,
			TZName:        v.TZName,
		}
		timezones = append(timezones, ct)
	}
	country := &Country{
		Name:           dto.Name,
		Iso3:           dto.Iso3,
		Iso2:           dto.Iso2,
		NumericCode:    dto.NumericCode,
		PhoneCode:      dto.PhoneCode,
		Capital:        dto.Capital,
		Currency:       dto.Currency,
		CurrencyName:   dto.CurrencyName,
		CurrencySymbol: dto.CurrencySymbol,
		TLD:            dto.TLD,
		Native:         dto.Native,
		Region:         region,
		Subregion:      Subregion,
		Nationality:    dto.Nationality,
		Timezones:      timezones,
		Translations:   dto.Translations,
		Location:       dto.Location,
		Emoji:          dto.Emoji,
		EmojiU:         dto.EmojiU,
	}
	result, e := db.Mongo.Collection(country.Coll()).InsertOne(context.TODO(), country)
	if e != nil {
		return nil, err.Mongo(e)
	}
	country.ID = result.InsertedID.(bson.ObjectID)
	return country, nil
}

// ================================================================
// states
// ================================================================

type State struct {
	ID          bson.ObjectID `bson:"_id,omitempty"          json:"id,omitempty"`
	Name        string        `bson:"name,omitempty"         json:"name,omitempty"`
	CountryID   bson.ObjectID `bson:"country_id,omitempty"   json:"countryId,omitempty"`
	CountryCode string        `bson:"country_code,omitempty" json:"countryCode,omitempty"`
	CountryName string        `bson:"country_name,omitempty" json:"countryName,omitempty"`
	Iso2        string        `bson:"iso2,omitempty"         json:"iso2,omitempty"`
	Iso3166_2   string        `bson:"iso3166_2,omitempty"    json:"iso3166_2,omitempty"`
	FipsCode    string        `bson:"fips_code,omitempty"    json:"fipsCode,omitempty"`
	Type        string        `bson:"type,omitempty"         json:"type,omitempty"`
	Level       int           `bson:"level,omitempty"        json:"level,omitempty"`
	ParentID    bson.ObjectID `bson:"parent_id,omitempty"    json:"parentId,omitempty"`
	Location    db.GeoPoint   `bson:"location"               json:"location"`
	Timezone    string        `bson:"timezone,omitempty"     json:"timezone,omitempty"`
	CreatedAt   time.Time     `bson:"created_at"             json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updated_at"             json:"updatedAt"`
	db.Odm      `bson:"-" json:"-"`
}

func NewState() *State {
	state := &State{}
	state.Odm.Model = state
	return state
}

func (s *State) CollectionName() string { return "states" }
func (s *State) GetID() bson.ObjectID   { return s.ID }
func (s *State) SetID(id bson.ObjectID) { s.ID = id }

func (s *State) BeforeCreate() error {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

func (s *State) BeforeUpdate() error {
	s.UpdatedAt = time.Now()
	return nil
}

// ================================================================
// cities
// ================================================================
type City struct {
	ID          bson.ObjectID `bson:"_id,omitempty"          json:"id,omitempty"`
	Name        string        `bson:"name,omitempty"         json:"name,omitempty"`
	StateID     bson.ObjectID `bson:"state_id,omitempty"     json:"stateId,omitempty"`
	StateCode   string        `bson:"state_code,omitempty"   json:"stateCode,omitempty"`
	StateName   string        `bson:"state_name,omitempty"   json:"stateName,omitempty"`
	CountryID   bson.ObjectID `bson:"country_id,omitempty"   json:"countryId,omitempty"`
	CountryCode string        `bson:"country_code,omitempty" json:"countryCode,omitempty"`
	CountryName string        `bson:"country_name,omitempty" json:"countryName,omitempty"`
	Location    db.GeoPoint   `bson:"location"               json:"location"`
	Timezone    string        `bson:"timezone,omitempty"     json:"timezone,omitempty"`
	WikiDataID  string        `bson:"wikiDataId,omitempty"   json:"wikiDataId,omitempty"`
	CreatedAt   time.Time     `bson:"created_at"             json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updated_at"             json:"updatedAt"`
	db.Odm      `bson:"-" json:"-"`
}

func NewCity() *City {
	city := &City{}
	city.Odm.Model = city
	return city
}

func (c *City) CollectionName() string { return "cities" }
func (c *City) GetID() bson.ObjectID   { return c.ID }
func (c *City) SetID(id bson.ObjectID) { c.ID = id }

func (c *City) BeforeCreate() error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *City) BeforeUpdate() error {
	c.UpdatedAt = time.Now()
	return nil
}
