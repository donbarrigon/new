package controller

import (
	"donbarrigon/new/internal/database/data/repository"
	"donbarrigon/new/internal/database/data/seed"
	"donbarrigon/new/internal/utils/handler"
	"donbarrigon/new/internal/utils/logs"
	"reflect"
	"runtime"
)

// ejecuta las semillas que no se han ejecutado
func Seed(c *handler.Context) {
	records := repository.GetSeedTracker()
	seeds := repository.GetRunSeeds()

	execute := []func(){}
	for _, s := range seeds {
		exists := false
		ptr := reflect.ValueOf(s).Pointer()
		name := runtime.FuncForPC(ptr).Name()

		for _, record := range records {
			if record["name"] == name {
				exists = true
				break
			}
		}
		if !exists {
			execute = append(execute, s)
		}
	}

	repository.RunSeeds(execute)
	logs.Info("✅ Semillas ejecutadas")
	c.ResponseNoContent()
}

// fuerza la ejecucion de todas las semillas
func SeedAll(c *handler.Context) {
	seeds := repository.GetRunSeeds()
	repository.RunSeeds(seeds)
	logs.Info("✅ Semillas ejecutadas")
	c.ResponseNoContent()
}

// fuerza la ejecucion de una semilla
func SeedRun(c *handler.Context) {
	seeds := repository.GetRunSeeds()
	execute := seed.Seeds{}
	searchName := c.Get("name")

	for _, s := range seeds {
		ptr := reflect.ValueOf(s).Pointer()
		name := runtime.FuncForPC(ptr).Name()
		if name == searchName {
			execute = append(execute, s)
		}
	}

	repository.RunSeeds(execute)
	logs.Info("✅ Semilla ejecutada")
	c.ResponseNoContent()
}
