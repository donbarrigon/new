package migration

import "donbarrigon/new/internal/utils/db"

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
