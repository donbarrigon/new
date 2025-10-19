package migration

import "donbarrigon/new/internal/utils/db"

func Run() []db.Migration {
	m := []db.Migration{
		CreateUsersCollection{},
		CreateRolesCollection{},
		CreatePremissionsCollection{},
		CreateTrashCollection{},
		CreateHistoryCollection{},
		CreateTokensCollection{},
		CreateSessionCollection{},
		CreateCountriesCollection{},
		CreateStatesCollection{},
		CreateCitiesCollection{},
	}
	return m
}
