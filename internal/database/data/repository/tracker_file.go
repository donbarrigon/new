package repository

import (
	"donbarrigon/new/internal/utils/logs"
	"os"
	"path/filepath"
)

const migrationFileName = "migration_tracker.txt"
const seedFileName = "seed_tracker.txt"
const tmpPath = "tmp"

func openFile(fileName string) *os.File {
	if e := os.MkdirAll(tmpPath, os.ModePerm); e != nil {
		logs.Error("Fail to create log directory %s: %s", fileName, e.Error())
		panic(e.Error())
	}

	filePath := filepath.Join(tmpPath, fileName)

	file, e := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if e != nil {
		logs.Error("Error al abrir el archivo: %s %s", filePath, e.Error())
		panic(e.Error())
	}
	return file
}

func dropFile(fileName string) {
	filePath := filepath.Join(tmpPath, fileName)
	if e := os.Remove(filePath); e != nil {
		panic(e.Error())
	}
}
