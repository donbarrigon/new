package schema

import "donbarrigon/new/internal/utils/db"

// ================================================================
// CreateCountriesCollection
// ================================================================

type CreateCountriesCollection struct{}

func (c CreateCountriesCollection) Name() string {
	return "countries"
}

func (c CreateCountriesCollection) Up() {
	db.CreateCollection(c.Name(), func(col *db.CollectionBuilder) {
		col.Index("t", "name")
	})
}

func (c CreateCountriesCollection) Down() {
	db.DropCollection(c.Name())
}

// ================================================================
// CreateStatesCollection
// ================================================================

type CreateStatesCollection struct{}

func (s CreateStatesCollection) Name() string {
	return "states"
}

func (s CreateStatesCollection) Up() {
	db.CreateCollection(s.Name(), func(col *db.CollectionBuilder) {
		col.Index("i", "country_id")
		col.Index("t", "name", "country_name")
	})
}

func (s CreateStatesCollection) Down() {
	db.DropCollection(s.Name())
}

// ================================================================
// CreateCitiesCollection
// ================================================================

type CreateCitiesCollection struct{}

func (c CreateCitiesCollection) Name() string {
	return "cities"
}

func (c CreateCitiesCollection) Up() {
	db.CreateCollection(c.Name(), func(col *db.CollectionBuilder) {
		col.Index("i", "name")
		col.Index("i", "state_id")
		col.Index("i", "country_id")
		col.Index("t", "name", "state_name", "country_name")
	})
}

func (c CreateCitiesCollection) Down() {
	db.DropCollection(c.Name())
}
