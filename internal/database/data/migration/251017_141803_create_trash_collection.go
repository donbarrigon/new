package migration

import "donbarrigon/new/internal/utils/db"

type CreateTrashCollection struct{}

func (u CreateTrashCollection) Name() string {
	return "trash"
}

func (u CreateTrashCollection) Up() {
	db.CreateCollection(u.Name(), func(col *db.CollectionBuilder) {
		col.Index("i", "user_id")
		col.Index("i", "collection")
		col.Index("i", "collection_id")
	})
}

func (u CreateTrashCollection) Down() {
	db.DropCollection(u.Name())
}
