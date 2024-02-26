package configs

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strconv"
)

type AppConfig struct {
	DBConfig *pgconn.Config
}

func ConfigInit() (*AppConfig, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}

	dbConfig, err := loadDBConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load DB config")
	}

	return &AppConfig{
		DBConfig: dbConfig,
	}, nil
}

func loadDBConfig() (*pgconn.Config, error) {
	user := viper.GetString("db.postgres.username")
	dbname := viper.GetString("db.postgres.dbname")
	host := viper.GetString("db.postgres.host")
	port := viper.GetString("db.postgres.port")
	password := viper.GetString("db.postgres.password")

	portUint, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse port")
	}

	return &pgconn.Config{
		Host:     host,
		Port:     uint16(portUint),
		User:     user,
		Password: password,
		Database: dbname,
	}, nil
}
