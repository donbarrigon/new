package controller

import (
	"bufio"
	"context"
	"donbarrigon/new/internal/database/data/migration"
	"donbarrigon/new/internal/database/handler/service"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/handler"
	"donbarrigon/new/internal/utils/logs"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func Migrate(c *handler.Context) {

	migrations := migration.Run()
	records := service.GetMigrationTracker()
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
	service.RunMigrations("up", execute)

	logs.Info("Migrations executed")
	c.NoContent()
}

func Rollback(ctx *handler.Context) {

	file := openFile("migration_tracker.txt")
	defer file.Close()

	migration.Migrations = []app.List{}
	migration.Run()

	scanner := bufio.NewScanner(file)
	records := []map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		record := map[string]string{}
		for _, field := range fields {
			parts := strings.SplitN(field, ":", 2)
			if len(parts) < 2 {
				continue
			}
			record[parts[0]] = parts[1]
		}
		if record["action"] == "up" {
			records = append(records, record)
		}
		if record["action"] == "down" {
			for i, recordUp := range records {
				if recordUp["name"] == record["name"] {
					records = append(records[:i], records[i+1:]...)
					break
				}
			}
		}
	}
	if er := scanner.Err(); er != nil {
		app.PrintError("Fail to read file: :file :error", app.E("file", file.Name()), app.E("error", er.Error()))
		panic(er.Error())
	}

	var last string
	if len(records) > 0 {
		last = records[len(records)-1]["executed_at"]
	}

	filtered := []map[string]string{}
	for _, record := range records {
		if record["executed_at"] == last {
			filtered = append(filtered, record)
		}
	}

	migrations := []app.List{}
	for _, filter := range filtered {
		for _, m := range migration.Migrations {
			if m.Get("name").(string) == filter["name"] {
				migrations = append(migrations, m)
			}
		}
	}
	runMigrations("down", migrations, file)

	app.PrintInfo("Migrations rolled back")
	ctx.ResponseNoContent()
}

func Reset(ctx *handler.Context) {

	file := openFile("migration_tracker.txt")
	defer file.Close()

	filePath := filepath.Join(app.Env.LOG_PATH, "seed_tracker.txt")
	er := os.Remove(filePath)
	if er != nil {
		if !os.IsNotExist(er) {
			fmt.Println("Fail to remove file seed_tracker:", er)
			return
		}
	}

	migration.Migrations = []app.List{}
	migration.Run()

	scanner := bufio.NewScanner(file)
	records := []map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		record := map[string]string{}
		for _, field := range fields {
			parts := strings.SplitN(field, ":", 2)
			if len(parts) < 2 {
				continue
			}
			record[parts[0]] = parts[1]
		}
		if record["action"] == "up" {
			records = append(records, record)
		}
		if record["action"] == "down" {
			for i, recordUp := range records {
				if recordUp["name"] == record["name"] {
					records = append(records[:i], records[i+1:]...)
					break
				}
			}
		}
	}
	if er := scanner.Err(); er != nil {
		app.PrintError("Fail to read file: :file :error", app.E("file", file.Name()), app.E("error", er.Error()))
		panic(er.Error())
	}

	migrations := []app.List{}
	for _, filter := range records {
		for _, m := range migration.Migrations {
			if m.Get("name").(string) == filter["name"] {
				migrations = append(migrations, m)
			}
		}
	}
	runMigrations("down", migrations, file)

	app.PrintInfo("Migration reset")
	ctx.ResponseNoContent()

}

func Refresh(ctx *handler.Context) {
	file := openFile("migration_tracker.txt")
	defer file.Close()

	filePath := filepath.Join(app.Env.LOG_PATH, "seed_tracker.txt")
	er := os.Remove(filePath)
	if er != nil {
		if !os.IsNotExist(er) {
			fmt.Println("Fail to remove file seed_tracker:", er)
			return
		}
	}

	migration.Migrations = []app.List{}
	migration.Run()

	scanner := bufio.NewScanner(file)
	records := []map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		record := map[string]string{}
		for _, field := range fields {
			parts := strings.SplitN(field, ":", 2)
			if len(parts) < 2 {
				continue
			}
			record[parts[0]] = parts[1]
		}
		if record["action"] == "up" {
			records = append(records, record)
		}
		if record["action"] == "down" {
			for i, recordUp := range records {
				if recordUp["name"] == record["name"] {
					records = append(records[:i], records[i+1:]...)
					break
				}
			}
		}
	}
	if er := scanner.Err(); er != nil {
		app.PrintError("Fail to read file: :file :error", app.E("file", file.Name()), app.E("error", er.Error()))
		panic(er.Error())
	}

	migrations := []app.List{}
	for _, filter := range records {
		for _, m := range migration.Migrations {
			if m.Get("name").(string) == filter["name"] {
				migrations = append(migrations, m)
			}
		}
	}
	runMigrations("down", migrations, file)
	runMigrations("up", migration.Migrations, file)

	app.PrintInfo("Migrations refreshed")
	ctx.ResponseNoContent()
}

func Fresh(ctx *handler.Context) {
	app.DB.Drop(context.TODO())

	filePath := filepath.Join(app.Env.LOG_PATH, "migration_tracker.txt")
	er := os.Remove(filePath)
	if er != nil {
		if !os.IsNotExist(er) {
			fmt.Println("Fail to remove file migration_tracker:", er)
			return
		}
	}

	filePath = filepath.Join(app.Env.LOG_PATH, "seed_tracker.txt")
	er = os.Remove(filePath)
	if er != nil {
		if !os.IsNotExist(er) {
			fmt.Println("Fail to remove file seed_tracker:", er)
			return
		}
	}

	file := openFile("migration_tracker.txt")
	defer file.Close()

	migration.Migrations = []app.List{}
	migration.Run()
	runMigrations("up", migration.Migrations, file)

	app.PrintInfo("Database refreshed")
	ctx.ResponseNoContent()
}
