package repository

import (
	"bufio"
	"context"
	"donbarrigon/new/internal/database/data/schema"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/logs"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func GetMigrationTracker() []map[string]string {
	file := openFile(migrationFileName)
	defer file.Close()
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
		records = append(records, record)
	}
	if e := scanner.Err(); e != nil {
		logs.Error("Fail to read file: %s %s", file.Name(), e.Error())
		panic(e.Error())
	}
	return records
}

func GetRunMigrations() []db.Migration {
	return schema.Run()
}

func GetMigrationRecordsUp() []map[string]string {
	records := GetMigrationTracker()
	recordsUp := []map[string]string{}

	for _, record := range records {
		if record["action"] == "up" {
			recordsUp = append(recordsUp, record)
		} else {
			for i, recordUp := range recordsUp {
				if recordUp["name"] == record["name"] {
					recordsUp = append(recordsUp[:i], recordsUp[i+1:]...)
					break
				}
			}
		}
	}

	return recordsUp
}

func SelectMigrationsByRecords(records []map[string]string) []db.Migration {
	migrations := schema.Run()
	result := []db.Migration{}
	for _, record := range records {
		for _, m := range migrations {
			name := reflect.TypeOf(m).Name()
			if record["name"] == name {
				result = append(result, m)
			}
		}
	}
	return result
}

func RunMigrations(action string, migrations []db.Migration) {
	file := openFile(migrationFileName)
	defer file.Close()
	executedAt := time.Now()
	for _, m := range migrations {
		if action == "up" {
			m.Up()
		} else if action == "down" {
			m.Down()
		} else {
			panic("Accion desconocida: " + action)
		}
		name := reflect.TypeOf(m).Name()
		line := fmt.Sprintf("executed_at:%s\taction:%s\tname:%s\n", executedAt, action, name)
		if _, e := file.WriteString(line); e != nil {
			logs.Error("Error al escribir en el archivo: %s %s", file.Name(), e.Error())
			panic(e.Error())
		}
		logs.Info(line)
	}
}

func DropDB() {
	db.DB.Drop(context.TODO())
}

func DropMigrationTracker() {
	dropFile(migrationFileName)
}
