package migration

import "donbarrigon/new/internal/utils/db"

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
