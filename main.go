package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host        string
	Port        int
	Dbname      string
	User        string
	Password    string
	Sslmode     string
	Sslrootcert string
}

type Config struct {
	Telegram_token string
	Database       DatabaseConfig
}

var config Config

var telegramUrl string

var DB *sql.DB = nil

func init() {
	fileData, err := ioutil.ReadFile(".config.json")
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(fileData, &config)
	if err != nil {
		log.Fatalln(err)
	}

	telegramUrl = fmt.Sprintf(
		"https://api.telegram.org/bot%s",
		config.Telegram_token,
	)

	postgresConn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s sslrootcert=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.Dbname,
		config.Database.User,
		config.Database.Password,
		config.Database.Sslmode,
		config.Database.Sslrootcert,
	)

	db, err := sql.Open("postgres", postgresConn)
	if err != nil {
		panic(err)
	}

	DB = db
}

func main() {
	defer DB.Close()
	TelegramPoll()
}
