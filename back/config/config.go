package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type Config_t struct {
	PostgresAdminUser     string `env:"POSTGRES_ADMIN_USER"`
	PostgresAdminPassword string `env:"POSTGRES_ADMIN_PASSWORD"`
	PostgresUser          string `env:"POSTGRES_USER"`
	PostgresPassword      string `env:"POSTGRES_PASSWORD"`
	PostgresDb            string `env:"POSTGRES_DB"`
	PostgresHost          string `env:"POSTGRES_HOST"`
	PostgresPort          string `env:"POSTGRES_PORT"`

	ISDebugMode string `env:"IS_DEBUG_MODE"`

	ServerPort string
}

var Config = Config_t{}

func LoadConfig() error {
	serverport := flag.String("port", "8080", "server port")
	flag.Parse()
	Config.ServerPort = *serverport

	err := godotenv.Load(".env")
	if err != nil {
		slog.Debug("Error loading .env file")
	}

	t := reflect.TypeOf(Config)
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		tag, ok := t.Field(i).Tag.Lookup("env")
		if !ok {
			continue
		}
		v, ok := os.LookupEnv(tag)
		if !ok {
			return fmt.Errorf("%s is not set", tag)
		}
		reflect.ValueOf(&Config).Elem().FieldByName(fieldName).SetString(v)
	}
	return nil
}
