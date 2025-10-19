package service

import (
	"bufio"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/logs"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

const migrationFileName = "migration_tracker.txt"
const tmpPath = "tmp"

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

func openFile(fileName string) *os.File {
	if e := os.MkdirAll(tmpPath, os.ModePerm); e != nil {
		logs.Error("Fail to create log directory %s: %s", fileName, e.Error())
		panic(e.Error())
	}

	filePath := filepath.Join(tmpPath, fileName)

	file, e := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if e != nil {
		logs.Error("Fail to open file: %s %s", filePath, e.Error())
		panic(e.Error())
	}
	return file
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
			panic("Unknown action: " + action)
		}
		name := reflect.TypeOf(m).Name()
		line := fmt.Sprintf("executed_at:%s\taction:%s\tname:%v\n", executedAt, action, name)
		if _, e := file.WriteString(line); e != nil {
			logs.Error("Fail to write %s %s", file.Name(), e.Error())
			panic(e.Error())
		}
		logs.Info(line)
	}
}
