package migration

import "donbarrigon/new/internal/utils/db"

type CreateRolesCollection struct{}

func (r CreateRolesCollection) Name() string {
	return "roles"
}

func (r CreateRolesCollection) Up() {
	db.CreateCollection(r.Name(), func(col *db.CollectionBuilder) {
		col.Index("unique", "name")
	})
}

func (r CreateRolesCollection) Down() {
	db.DropCollection(r.Name())
}
