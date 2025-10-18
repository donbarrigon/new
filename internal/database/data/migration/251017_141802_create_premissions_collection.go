package migration

import "donbarrigon/new/internal/utils/db"

type CreatePremissionsCollection struct {
}

func (u CreatePremissionsCollection) Name() string {
	return "premissions"
}

func (u CreatePremissionsCollection) Up() {
	db.CreateCollection(u.Name(), func(col *db.CollectionBuilder) {
		col.Index("unique", "name")
	})
}

func (u CreatePremissionsCollection) Down() {
	db.DropCollection(u.Name())
}
