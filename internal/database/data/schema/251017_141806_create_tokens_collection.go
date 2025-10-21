package schema

import "donbarrigon/new/internal/utils/db"

type CreateTokensCollection struct{}

func (t CreateTokensCollection) Name() string {
	return "tokens"
}

func (t CreateTokensCollection) Up() {
	db.CreateCollection(t.Name(), func(col *db.CollectionBuilder) {
		col.Index("i", "user_id")
		col.Index("u", "token")
	})
}

func (t CreateTokensCollection) Down() {
	db.DropCollection(t.Name())
}
