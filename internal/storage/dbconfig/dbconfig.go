package dbconfig

import (
	"fmt"
	"github.com/spf13/viper"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Name     string
	Port     int
}

func NewDBConfig() (*DBConfig, error) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Читаем конфигурацию
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var dbConfig DBConfig
	err = viper.UnmarshalKey("db", &dbConfig)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &dbConfig, nil

}
