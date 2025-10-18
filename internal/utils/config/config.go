package config

import (
	"bufio"
	"donbarrigon/new/internal/utils/logs"
	"os"
	"strconv"
	"strings"
)

type V map[string]any

// si se proporciona una ruta, se usa la primera; de lo contrario, se carga el archivo .env por defecto
var (
	AppName   string = "NewApp"
	AppKey    string = ""
	AppURL    string = "http://localhost:3000"
	AppLocale string = "es"
	AppDebug  bool   = false

	ServerPort          string = "3000"
	ServerHttpsEnabled  bool   = false
	ServerHttpsCertPath string = "./certs"
	ServerHostWhiteList string = ""
	ServerReadTimeout   int    = 30
	ServerWriteTimeout  int    = 30

	SessionLifetime int = 10080
	SessionDriver   int = SESSION_DRIVER_FILE

	DbName             string = "samplemflix"
	DbConnectionString string = "mongodb://localhost:27017"
	DbMigrationEnable  bool   = false

	MailHost     string = "smtp.gmail.com"
	MailPort     string = "587"
	MailUsername string = "email@gmail.com"
	MailPassword string = "secreto123"
	MailFromName string = "Don Barrigon"
	MailIdentity string = "donbarrigon@gmail.com"
)

const SESSION_DRIVER_FILE = 0
const SESSION_DRIVER_MONGO = 1
const SESSION_DRIVER_REDIS = 2

func LoadEnv(filepath ...string) {
	f := ".env"
	if len(filepath) > 0 {
		f = filepath[0]
	}

	file, err := os.Open(f)
	if err != nil {
		logs.Error("no fue posible abrir el archivo %v: %v", f, err.Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0

	for scanner.Scan() {
		i++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			logs.Warning("Error de formato en la línea %v: %v", i, line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if idx := strings.Index(value, "#"); idx != -1 && !strings.HasPrefix(value, `"`) && !strings.HasPrefix(value, `'`) {
			value = strings.TrimSpace(value[:idx])
		}
		value = strings.Trim(value, `"'`)

		if key == "" {
			logs.Warning("Clave vacía detectada al cargar variables de entorno en la línea %v: %v", i, line)
			continue
		}

		os.Setenv(key, value)

		switch key {
		case "APP_NAME":
			AppName = value
		case "APP_KEY":
			AppKey = value
		case "APP_URL":
			AppURL = value
		case "APP_LOCALE":
			AppLocale = value

		case "SERVER_PORT":
			ServerPort = value
		case "SERVER_HTTPS_ENABLED":
			ServerHttpsEnabled = false
			if strings.ToLower(value) == "true" {
				ServerHttpsEnabled = true
			}
		case "SERVER_HTTPS_CERT_PATH":
			ServerHttpsCertPath = value
		case "SERVER_HOST_WHITE_LIST":
			ServerHostWhiteList = value
		case "SERVER_READ_TIMEOUT":
			timeout, e := strconv.Atoi(value)
			if e != nil {
				ServerReadTimeout = timeout
			}
		case "SERVER_WRITE_TIMEOUT":
			timeout, e := strconv.Atoi(value)
			if e != nil {
				ServerWriteTimeout = timeout
			}
		case "SESSION_LIFETIME":
			duration, e := strconv.Atoi(value)
			if e != nil {
				SessionLifetime = duration
			}
		case "DB_MIGRATION_ENABLE":
			DbMigrationEnable = false
			if strings.ToLower(value) == "true" {
				DbMigrationEnable = true
			}

		case "DB_NAME":
			DbName = value
		case "DB_CONNECTION_STRING":
			DbConnectionString = value

		case "LOG_LEVEL":
			switch strings.ToUpper(strings.TrimSpace(value)) {
			case "OFF":
				logs.LV = logs.LV_OFF
			case "EMERGENCY":
				logs.LV = logs.LV_EMERGENCY
			case "ALERT":
				logs.LV = logs.LV_ALERT
			case "CRITICAL":
				logs.LV = logs.LV_CRITICAL
			case "ERROR":
				logs.LV = logs.LV_ERROR
			case "WARNING":
				logs.LV = logs.LV_WARNING
			case "NOTICE":
				logs.LV = logs.LV_NOTICE
			case "INFO":
				logs.LV = logs.LV_INFO
			case "DEBUG":
				logs.LV = logs.LV_DEBUG
			default:
				logs.LV = logs.LV_DEBUG
			}
		case "LOG_FLAGS":
			flags := 0
			parts := strings.Split(value, ",")
			for _, part := range parts {
				switch strings.ToUpper(strings.TrimSpace(part)) {
				case "TIMESTAMP":
					flags |= logs.FLAG_TIMESTAMP
				case "FILE":
					flags |= logs.FLAG_FILE
				case "SHORTFILE":
					flags |= logs.FLAG_SHORTFILE
				case "LEVEL":
					flags |= logs.FLAG_LEVEL
				case "CONSOLE_AS_JSON":
					flags |= logs.FLAG_CONSOLE_AS_JSON
				}
			}
			logs.Flags = flags
		case "LOG_OUTPUT":
			outputs := 0
			parts := strings.Split(value, ",")
			for _, part := range parts {
				switch strings.ToUpper(strings.TrimSpace(part)) {
				case "CONSOLE":
					outputs |= logs.OUTPUT_CONSOLE
				case "FILE":
					outputs |= logs.OUTPUT_FILE
				case "REMOTE":
					outputs |= logs.OUTPUT_REMOTE
				}
			}
			logs.Outputs = outputs
		case "LOG_PATH":
			logs.Path = value
		case "LOG_CHANNEL":
			value = strings.ToLower(value)
			if value == "monthly" || value == "weekly" || value == "single" {
				logs.Channel = value
			} else {
				logs.Channel = "daily"
			}
		case "LOG_DAYS":
			days, err := strconv.Atoi(value)
			if err != nil {
				logs.Warning("LOG_DAYS valor inválido en la línea %v: %v", i, value)
				continue
			}
			logs.Days = days
		case "LOG_DATE_FORMAT":
			logs.DateFormat = value
		case "MAIL_HOST":
			MailHost = value
		case "MAIL_PORT":
			MailPort = value
		case "MAIL_USERNAME":
			MailUsername = value
		case "MAIL_PASSWORD":
			MailPassword = value
		case "MAIL_FROM_NAME":
			MailFromName = value
		case "MAIL_IDENTITY":
			MailIdentity = value
		default:
			logs.Warning("La variable de entorno %v no existe", key)
		}
	}

	if scanner.Err() != nil {
		logs.Error("Error al leer el archivo %v: %v", f, scanner.Err().Error())
		return
	}

	logs.Info("Configuraciones cargadas: %v", f)
}
