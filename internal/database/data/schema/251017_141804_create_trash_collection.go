package schema

import "donbarrigon/new/internal/utils/db"

type CreateTrashCollection struct{}

func (t CreateTrashCollection) Name() string {
	return "trash"
}

func (t CreateTrashCollection) Up() {
	db.CreateCollection(t.Name(), func(col *db.CollectionBuilder) {
		col.Index("i", "user_id")
		col.Index("i", "collection")
		col.Index("i", "collection_id")
	})
}

func (t CreateTrashCollection) Down() {
	db.DropCollection(t.Name())
}
