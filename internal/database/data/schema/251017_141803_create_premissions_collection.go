package schema

import "donbarrigon/new/internal/utils/db"

type CreatePremissionsCollection struct {
}

func (p CreatePremissionsCollection) Name() string {
	return "premissions"
}

func (p CreatePremissionsCollection) Up() {
	db.CreateCollection(p.Name(), func(col *db.CollectionBuilder) {
		col.Index("unique", "name")
	})
}

func (p CreatePremissionsCollection) Down() {
	db.DropCollection(p.Name())
}
