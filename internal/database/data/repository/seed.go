package repository

import (
	"bufio"
	"donbarrigon/new/internal/database/data/seed"
	"donbarrigon/new/internal/utils/logs"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func GetSeedTracker() []map[string]string {
	file := openFile(seedFileName)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	records := []map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		record := map[string]string{}
		for _, field := range fields {
			parts := strings.SplitN(field, ":", 2)
			if len(parts) < 1 {
				continue
			}
			record[parts[0]] = parts[1]
		}
		records = append(records, record)
	}
	if e := scanner.Err(); e != nil {
		logs.Error("Error al leer el archivo de las semillas: %s %s", file.Name(), e.Error())
		panic(e.Error())
	}
	return records
}

func GetRunSeeds() seed.Seeds {
	return seed.Run()
}

func RunSeeds(seeds seed.Seeds) {
	file := openFile(migrationFileName)
	defer file.Close()
	executedAt := time.Now()
	for _, s := range seeds {

		s()
		ptr := reflect.ValueOf(s).Pointer()
		name := runtime.FuncForPC(ptr).Name()

		line := fmt.Sprintf("executed_at:%s\tname:%s\n", executedAt, name)
		if _, e := file.WriteString(line); e != nil {
			logs.Error("Error al escribir en el archivo: %s %s", file.Name(), e.Error())
			panic(e.Error())
		}
		logs.Info(line)
	}
}

func DropSeedTracker() {
	dropFile(seedFileName)
}
