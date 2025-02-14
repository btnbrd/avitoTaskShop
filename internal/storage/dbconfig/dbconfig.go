package dbconfig

import (
	"github.com/lpernett/godotenv"
	"os"
	"strconv"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Name     string
	Port     int
}

//func NewDBConfig() (*DBConfig, error) {
//	return &DBConfig{
//		Host:     os.Getenv("DB_HOST"),
//		Port:     5432,
//		User:     os.Getenv("DB_USER"),
//		Password: os.Getenv("DB_PASSWORD"),
//		Name:     os.Getenv("DB_NAME"),
//	}, nil
//}

func NewDBConfig() (*DBConfig, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	portString := os.Getenv("DB_PORT")
	port, _ := strconv.Atoi(portString)

	return &DBConfig{Name: name, User: user, Password: password, Host: host, Port: port}, nil

}
