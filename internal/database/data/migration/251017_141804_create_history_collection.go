package migration

import "donbarrigon/new/internal/utils/db"

type CreateHistoryCollection struct{}

func (u CreateHistoryCollection) Name() string {
	return "history"
}

func (u CreateHistoryCollection) Up() {
	db.CreateCollection(u.Name(), func(col *db.CollectionBuilder) {
		col.Index("i", "user_id")
		col.Index("i", "collection")
	})
}

func (u CreateHistoryCollection) Down() {
	db.DropCollection(u.Name())
}
