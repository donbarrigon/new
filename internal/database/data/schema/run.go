package schema

import "donbarrigon/new/internal/utils/db"

func Run() []db.Migration {
	return []db.Migration{
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
}
