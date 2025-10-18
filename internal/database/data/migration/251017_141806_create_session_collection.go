package migration

import "donbarrigon/new/internal/utils/db"

type CreateSessionCollection struct {
	db.Migration
}

func (m CreateSessionCollection) Name() string {
	return "session"
}

func (m CreateSessionCollection) Up() {
	db.CreateCollection(m.Name(), func(col *db.CollectionBuilder) {
		col.Index("i", "user_id")
		col.Index("u", "token")
	})
}

func (m CreateSessionCollection) Down() {
	db.DropCollection(m.Name())
}
