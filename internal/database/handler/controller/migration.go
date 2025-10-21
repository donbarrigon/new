package controller

import (
	"donbarrigon/new/internal/database/data/repository"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/handler"
	"donbarrigon/new/internal/utils/logs"
	"reflect"
)

// ejecuta up() de las migraciones sin ejecutar
func Migrate(c *handler.Context) {

	migrations := repository.GetRunMigrations()
	records := repository.GetMigrationTracker()
	execute := []db.Migration{}

	for _, m := range migrations {
		exists := false
		t := reflect.TypeOf(m)
		name := t.Name()
		for _, record := range records {
			if record["name"] == name {
				if record["action"] == "up" {
					exists = true
				}
				if record["action"] == "down" {
					exists = false
				}
			}
		}
		if !exists {
			execute = append(execute, m)
		}
	}
	repository.RunMigrations("up", execute)

	logs.Info("✅ Migraciones ejecutadas")
	c.ResponseNoContent()
}

// ejecuta down() de la ultima tanda
func Rollback(c *handler.Context) {

	recordsUp := repository.GetMigrationRecordsUp()

	var last string
	if len(recordsUp) > 0 {
		last = recordsUp[len(recordsUp)-1]["executed_at"]
	}

	filtered := []map[string]string{}
	for _, record := range recordsUp {
		if record["executed_at"] == last {
			filtered = append(filtered, record)
		}
	}

	execute := repository.SelectMigrationsByRecords(filtered)
	repository.RunMigrations("down", execute)

	logs.Info("✅ Migrations rolled back")
	c.ResponseNoContent()
}

// ejecuta down() de todas las migraciones que se hallan ejecutado al momento.
func Reset(c *handler.Context) {

	recordsUp := repository.GetMigrationRecordsUp()
	execute := repository.SelectMigrationsByRecords(recordsUp)
	repository.RunMigrations("down", execute)
	repository.DropSeedTracker()

	logs.Info("✅ Migraciones reseteadas")
	c.ResponseNoContent()
}

// hace reset y luego migrate (down + up).
func Refresh(c *handler.Context) {

	recordsUp := repository.GetMigrationRecordsUp()
	execute := repository.SelectMigrationsByRecords(recordsUp)
	repository.RunMigrations("down", execute)
	repository.DropSeedTracker()

	logs.Info("✅ Migraciones reseteadas")
	Migrate(c)
}

// elimina la base de datos y los trakers de las migraciones y los seed
// ejecuta migrate
func Fresh(c *handler.Context) {

	repository.DropDB()
	repository.DropMigrationTracker()
	repository.DropSeedTracker()

	logs.Info("✅ Migraciones refrescadas")
	Migrate(c)
}
