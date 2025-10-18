package migration

import "donbarrigon/new/internal/utils/db"

type CreateTokensCollection struct{}

func (u CreateTokensCollection) Name() string {
	return "tokens"
}

func (u CreateTokensCollection) Up() {
	db.CreateCollection(u.Name(), func(col *db.CollectionBuilder) {
		col.Index("i", "user_id")
		col.Index("u", "token")
	})
}

func (u CreateTokensCollection) Down() {
	db.DropCollection(u.Name())
}
