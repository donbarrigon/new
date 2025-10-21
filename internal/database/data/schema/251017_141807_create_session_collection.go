package schema

import "donbarrigon/new/internal/utils/db"

type CreateSessionCollection struct {
	db.Migration
}

func (s CreateSessionCollection) Name() string {
	return "session"
}

func (s CreateSessionCollection) Up() {
	db.CreateCollection(s.Name(), func(col *db.CollectionBuilder) {
		col.Index("i", "user_id")
		col.Index("u", "token")
	})
}

func (s CreateSessionCollection) Down() {
	db.DropCollection(s.Name())
}
