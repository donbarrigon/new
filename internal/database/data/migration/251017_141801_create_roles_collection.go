package migration

import "donbarrigon/new/internal/utils/db"

type CreateRolesCollection struct{}

func (u CreateRolesCollection) Name() string {
	return "roles"
}

func (u CreateRolesCollection) Up() {
	db.CreateCollection(u.Name(), func(col *db.CollectionBuilder) {
		col.Index("unique", "name")
	})
}

func (u CreateRolesCollection) Down() {
	db.DropCollection(u.Name())
}
