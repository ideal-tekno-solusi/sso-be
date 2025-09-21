package bootstrap

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

type Container struct {
	dbr *pgx.Conn
	dbw *pgx.Conn
}

func InitContainer() *Container {
	readEnv()
	viper.SetConfigName(os.Getenv("CONFIG_FILE"))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to read config with error: %v", err))
	}

	return &Container{}
}

func readEnv() {
	err := gotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to read env with error: %v", err))
	}
}
