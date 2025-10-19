package migration

import "donbarrigon/new/internal/utils/db"

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
