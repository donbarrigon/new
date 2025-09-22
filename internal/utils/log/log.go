package log

import (
	"donbarrigon/new/internal/utils/fm"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type LogLevel int
type LogFileFormat int

type Logger struct {
	ID           string   `json:"id,omitempty"      yaml:"id,omitempty"      xml:"id,omitempty"`
	Time         string   `json:"time,omitempty"    yaml:"time,omitempty"    xml:"time,omitempty"`
	Level        LogLevel `json:"level,omitempty"   yaml:"level,omitempty"   xml:"level,omitempty"`
	Message      string   `json:"message"           yaml:"message"           xml:"message"`
	Line         string   `json:"line,omitempty"    yaml:"line,omitempty"    xml:"line,omitempty"`
	File         string   `json:"file,omitempty"    yaml:"file,omitempty"    xml:"file,omitempty"`
	Placeholders []any    `json:"context,omitempty" yaml:"context,omitempty" xml:"context,omitempty"`
}

const (
	OFF       LogLevel = iota // 0 - Desactiva todos los logs
	EMERGENCY                 // 1 - El sistema está inutilizable
	ALERT                     // 2 - Se necesita acción inmediata
	CRITICAL                  // 3 - Fallo crítico del sistema
	ERROR                     // 4 - Errores de ejecución
	WARNING                   // 5 - Algo inesperado pasó
	NOTICE                    // 6 - Eventos normales, pero significativos
	INFO                      // 7 - Información general
	DEBUG                     // 8 - Información detallada para depuración
	PRINT                     // 9 - Solo imprime en consola
)

const (
	FLAG_TIMESTAMP       = 1 << iota // 1    - Agrega la fecha y hora formateada según DATE_FORMAT
	FLAG_FILE                        // 2    - Agrega nombre del archivo y número de línea: /a/b/c/d.go:23
	FLAG_SHORTFILE                   // 4    - Quita la ruta del nombre del archivo: d.go:23
	FLAG_LEVEL                       // 8    - Agrega el lv anteslog.LV = log. del mensaje (por ejemplo: [DEBUG])
	FLAG_CONSOLE_AS_JSON             // 16   - Salida en formato JSON en la consola
	FLAG_CONSOLE_COLOR               // 32   - Salida en consola con solor segun el lv
)

const (
	OUTPUT_CONSOLE = 1 << iota // 1 - salida por consola estándar
	OUTPUT_FILE                // 2 - salida a archivo
	OUTPUT_REMOTE              // 4 - enviar a un servidor remoto (opcional)
)

const (
	FILE_FORMAT_NDJSON LogFileFormat = iota // 0 - NDJSON (JSON por línea)
	FILE_FORMAT_CSV                         // 1 - CSV (valores separados por coma)
	FILE_FORMAT_PLAIN                       // 2 - Texto plano
	FILE_FORMAT_LTSV                        // 3 - LTSV (Labelled Tab Separated Values)
)

var LV LogLevel = DEBUG
var Flags = FLAG_TIMESTAMP | FLAG_FILE | FLAG_SHORTFILE | FLAG_LEVEL | FLAG_CONSOLE_AS_JSON | FLAG_CONSOLE_COLOR
var Outputs = OUTPUT_CONSOLE | OUTPUT_FILE
var FileFormat = FILE_FORMAT_NDJSON
var DateFormat = "2006-01-02 15:04:05.000"
var Path = "./log"
var Channel = "daily"
var Days = 30

func (lv LogLevel) String() string {
	switch lv {
	case OFF:
		return "OFF"
	case EMERGENCY:
		return "EMERGENCY"
	case ALERT:
		return "ALERT"
	case CRITICAL:
		return "CRITICAL"
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARNING"
	case NOTICE:
		return "NOTICE"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case PRINT:
		return "PRINT"
	default:
		return "UNKNOWN"
	}
}

func (lv LogLevel) Color() string {
	switch lv {
	case EMERGENCY:
		return "\033[91m" // rojo brillante
	case ALERT:
		return "\033[95m" // magenta
	case CRITICAL:
		return "\033[35m" // fucsia
	case ERROR:
		return "\033[31m" // rojo
	case WARNING:
		return "\033[33m" // amarillo
	case NOTICE:
		return "\033[92m" // verde claro
	case INFO:
		return "\033[34m" // azul
	case DEBUG:
		return "\033[36m" // cian
	case PRINT:
		return "\033[90m" // gris claro
	default:
		return "\033[0m"
	}
}

func (f LogFileFormat) String() string {
	switch f {
	case FILE_FORMAT_NDJSON:
		return "ndjson"
	case FILE_FORMAT_CSV:
		return "csv"
	case FILE_FORMAT_PLAIN:
		return "plain"
	case FILE_FORMAT_LTSV:
		return "ltsv"
	default:
		return "unknown"
	}
}

func (lv LogLevel) DefaultColor() string {
	return "\033[0m"
}

func (l LogLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

func (l LogLevel) MarshalYAML() (interface{}, error) {
	return l.String(), nil
}

func (l LogLevel) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(l.String(), start)
}

func Emergency(msg string, ph ...any) {
	if LV >= EMERGENCY {
		l := &Logger{
			Level:        EMERGENCY,
			Message:      msg,
			Placeholders: ph,
		}
		l.output()
	}
}

func Alert(msg string, ph ...any) {
	if LV >= ALERT {
		l := &Logger{
			Level:        ALERT,
			Message:      msg,
			Placeholders: ph,
		}
		l.output()
	}
}

func Critical(msg string, ph ...any) {
	if LV >= CRITICAL {
		l := &Logger{
			Level:        CRITICAL,
			Message:      msg,
			Placeholders: ph,
		}
		l.output()
	}
}

func Error(msg string, ph ...any) {
	if LV >= ERROR {
		l := &Logger{
			Level:        ERROR,
			Message:      msg,
			Placeholders: ph,
		}
		l.output()
	}
}

func Warning(msg string, ph ...any) {
	if LV >= WARNING {
		l := &Logger{
			Level:        WARNING,
			Message:      msg,
			Placeholders: ph,
		}
		l.output()
	}
}

func Notice(msg string, ph ...any) {
	if LV >= NOTICE {
		l := &Logger{
			Level:        NOTICE,
			Message:      msg,
			Placeholders: ph,
		}
		l.output()
	}
}

func Info(msg string, ph ...any) {
	if LV >= INFO {
		l := &Logger{
			Level:        INFO,
			Message:      msg,
			Placeholders: ph,
		}
		l.output()
	}
}

func Debug(msg string, ph ...any) {
	if LV >= DEBUG {
		l := &Logger{
			Level:        DEBUG,
			Message:      msg,
			Placeholders: ph,
		}
		l.output()
	}
}

func Log(level LogLevel, msg string, ph ...any) {
	if LV >= level {
		l := &Logger{
			Level:        level,
			Message:      msg,
			Placeholders: ph,
		}
		l.output()
	}
}

func Print(msg string, ph ...any) {
	l := &Logger{
		Level:        PRINT,
		Message:      msg,
		Placeholders: ph,
	}
	l.output()
}

func Dump(a any) {
	fmt.Println(formatDump(a))
}

func (l *Logger) DumpMany(vars ...any) {
	sep := strings.Repeat("-", 30)

	for i, v := range vars {
		if i > 0 {
			fmt.Println(sep)
		}
		fmt.Println(formatDump(v))
	}
}

func (l *Logger) output() {
	// Obtener información del runtime
	_, file, line, _ := runtime.Caller(2)

	// Preparar mensaje
	l.Message = fmt.Sprintf(l.Message, l.Placeholders...)

	if Flags&FLAG_TIMESTAMP != 0 {
		now := time.Now().Format(DateFormat)
		l.Time = now
	}

	if Flags&(FLAG_FILE) != 0 {
		l.File = file
		l.Line = strconv.Itoa(line)
	}

	if Flags&FLAG_SHORTFILE != 0 {
		l.File = filepath.Base(file)
		l.Line = strconv.Itoa(line)
	}

	if Outputs&OUTPUT_CONSOLE != 0 || l.Level == PRINT {
		l.outputConsole()
		if l.Level == PRINT {
			return
		}
	}

	if Outputs&OUTPUT_FILE != 0 {
		l.outputFile()
	}

	if Outputs&OUTPUT_REMOTE != 0 {
		l.outputRemote()
	}

}

func (l *Logger) outputConsole() {
	if Flags&FLAG_CONSOLE_AS_JSON != 0 {
		data, _ := json.MarshalIndent(l, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println(l.outputPlain(true))
	}
}

func (l *Logger) outputFile() {
	file := l.openFile()
	if file == nil {
		return
	}
	defer file.Close()

	l.deleteOldFiles()

	var output string

	switch FileFormat {
	case FILE_FORMAT_NDJSON:
		output = l.outputNDJSON()
	case FILE_FORMAT_CSV:
		output = l.outputCSV() // CSV: Time, Level, Message, Function, File, Line, context
	case FILE_FORMAT_PLAIN:
		output = l.outputPlain(false)
	case FILE_FORMAT_LTSV:
		output = l.outputLTSV()
	default:
		// Fallback a ndjson
		output = l.outputNDJSON()
	}

	file.WriteString(output + "\n")
}

func (l *Logger) outputRemote() {

}

func (l *Logger) outputPlain(withColor bool) string {
	var b strings.Builder
	color := ""
	reset := ""
	if Flags&FLAG_CONSOLE_COLOR != 0 && withColor {
		color = l.Level.Color()
		reset = l.Level.DefaultColor()
	}

	if Flags&FLAG_TIMESTAMP != 0 {
		b.WriteString(fmt.Sprintf("%s ", l.Time))
	}
	if Flags&FLAG_LEVEL != 0 {
		b.WriteString(fmt.Sprintf("[%s%s%s] ", color, l.Level.String(), reset))
	}

	if Flags&FLAG_CONSOLE_COLOR != 0 {
		b.WriteString(color + l.Message + reset)
	} else {
		b.WriteString(l.Message)
	}

	if Flags&(FLAG_FILE|FLAG_SHORTFILE) != 0 {
		b.WriteString(fmt.Sprintf(" (%s:%s)", l.File, l.Line))
	}

	return b.String()
}

func (l *Logger) outputNDJSON() string {
	jsonData, err := json.Marshal(l)
	var output string
	if err != nil {
		msg := "Log serialization error: " + err.Error()
		escapedDump := strings.ReplaceAll(formatDump(l), `"`, `\"`)
		escapedDump = strings.ReplaceAll(escapedDump, "\n", " ")
		escapedDump = strings.ReplaceAll(escapedDump, "\r", " ")
		ph := fm.Placeholder{"msg": msg, "log": escapedDump}
		output = ph.Replace(`{"level":"ERROR","message":":msg","log":":log"}`)

		Print(msg)
	} else {
		output = string(jsonData)
	}
	return output
}

func (l *Logger) outputCSV() string {
	var record []string

	if Flags&FLAG_TIMESTAMP != 0 {
		record = append(record, l.Time)
	}

	if Flags&FLAG_LEVEL != 0 {
		record = append(record, l.Level.String())
	}

	// El mensaje siempre va
	record = append(record, l.Message)

	if Flags&(FLAG_FILE|FLAG_SHORTFILE) != 0 {
		record = append(record, l.File)
		record = append(record, l.Line)
	}

	var b strings.Builder
	writer := csv.NewWriter(&b)
	writer.Write(record)
	writer.Flush()

	return strings.TrimSpace(b.String())
}

func (l *Logger) outputLTSV() string {
	escape := func(s string) string {
		s = strings.ReplaceAll(s, "\t", " ")
		s = strings.ReplaceAll(s, "\n", " ")
		s = strings.ReplaceAll(s, "\r", " ")
		return s
	}

	var b strings.Builder

	if Flags&FLAG_TIMESTAMP != 0 {
		b.WriteString("time:" + escape(l.Time) + "\t")
	}

	b.WriteString("level:" + escape(l.Level.String()) + "\t")
	b.WriteString("message:" + escape(l.Message) + "\t")

	if Flags&(FLAG_FILE|FLAG_SHORTFILE) != 0 {
		b.WriteString("file:" + escape(l.File) + "\t")
		b.WriteString("line:" + escape(l.Line) + "\t")
	}

	// Eliminar el tab final si existe
	output := b.String()
	if len(output) > 0 && output[len(output)-1] == '\t' {
		output = output[:len(output)-1]
	}

	return output
}

func (l *Logger) openFile() *os.File {
	var filename string
	now := time.Now()

	switch strings.ToLower(Channel) {
	case "daily":
		filename = now.Format("2006-01-02") + ".log"
	case "monthly", "mensual":
		filename = now.Format("2006-01") + ".log"
	case "weekly":
		year, week := now.ISOWeek()
		filename = fmt.Sprintf("%d-W%02d.log", year, week)
	default:
		filename = "output.log"
	}

	if err := os.MkdirAll(Path, os.ModePerm); err != nil {
		Print("No se pudo crear el directorio de logs: %v\n", err.Error())
		return nil
	}

	filePath := filepath.Join(Path, filename)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Print("Failed to create log directory: %v\n", err.Error())
		return nil
	}
	return file
}

func (l *Logger) deleteOldFiles() {
	if Channel != "single" {
		now := time.Now()
		// Eliminar archivos
		if Channel == "daily" && Days > 0 {
			entries, _ := os.ReadDir(Path)
			cutoff := now.AddDate(0, 0, -Days)
			for _, Entry := range entries {
				if Entry.IsDir() {
					continue
				}
				name := Entry.Name()
				if !strings.HasSuffix(name, ".log") {
					continue
				}
				datePart := strings.TrimSuffix(name, ".log")
				EntryDate, err := time.Parse("2006-01-02", datePart)
				if err == nil && EntryDate.Before(cutoff) {
					_ = os.Remove(filepath.Join(Path, name))
				}
			}
		}

		if Channel == "weekly" && Days > 0 {
			entries, _ := os.ReadDir(Path)

			// Calcular semanas a conservar, redondeando hacia arriba (mínimo 1)
			weeksToKeep := (Days + 6) / 7
			if weeksToKeep < 1 {
				weeksToKeep = 1
			}

			// Crear Objecta de semanas válidas (formato YYYY-Www)
			validWeeks := make(map[string]bool)
			for i := 0; i < weeksToKeep; i++ {
				weekTime := now.AddDate(0, 0, -7*i)
				year, week := weekTime.ISOWeek()
				weekStr := fmt.Sprintf("%d-W%02d", year, week)
				validWeeks[weekStr] = true
			}

			// Eliminar logs fuera del rango de semanas válidas
			for _, Entry := range entries {
				if Entry.IsDir() || !strings.HasSuffix(Entry.Name(), ".log") {
					continue
				}
				name := strings.TrimSuffix(Entry.Name(), ".log")

				// Formato semanal esperado: YYYY-Wxx
				if strings.Count(name, "-") == 1 && strings.Contains(name, "W") && len(name) == 8 {
					if !validWeeks[name] {
						_ = os.Remove(filepath.Join(Path, Entry.Name()))
					}
				}
			}
		}

		if strings.ToLower(Channel) == "monthly" && Days > 0 {
			entries, _ := os.ReadDir(Path)

			// Redondear hacia arriba los días a meses (mínimo 1)
			monthsToKeep := (Days + 29) / 30
			if monthsToKeep < 1 {
				monthsToKeep = 1
			}

			// Generar meses válidos
			validMonths := make(map[string]bool)
			for i := 0; i < monthsToKeep; i++ {
				month := now.AddDate(0, -i, 0).Format("2006-01")
				validMonths[month] = true
			}

			// Eliminar archivos fuera del rango permitido
			for _, Entry := range entries {
				if Entry.IsDir() || !strings.HasSuffix(Entry.Name(), ".log") {
					continue
				}
				name := strings.TrimSuffix(Entry.Name(), ".log")

				// Formato YYYY-MM
				if len(name) == 7 && strings.Count(name, "-") == 1 {
					if !validMonths[name] {
						_ = os.Remove(filepath.Join(Path, Entry.Name()))
					}
				}
			}
		}
	}
}
