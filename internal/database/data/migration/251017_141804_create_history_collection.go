package migration

import "donbarrigon/new/internal/utils/db"

type CreateHistoryCollection struct{}

func (h CreateHistoryCollection) Name() string {
	return "history"
}

func (h CreateHistoryCollection) Up() {
	db.CreateCollection(h.Name(), func(col *db.CollectionBuilder) {
		col.Index("i", "user_id")
		col.Index("i", "collection")
	})
}

func (h CreateHistoryCollection) Down() {
	db.DropCollection(h.Name())
}
