package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/mcherdakov/telegoat"
)

type DatabaseConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	DBName      string `json:"db_name"`
	User        string `json:"user"`
	Password    string `json:"password"`
	SSLMode     string `json:"ssl_mode"`
	SSLRootCert string `json:"ssl_root_cert"`
}

type Config struct {
	TelegramToken string `json:"telegram_token"`
	Database      DatabaseConfig
}

var config Config

var telegramClient telegoat.TelegramClient

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

	telegramClient = telegoat.NewTelegramClient(config.TelegramToken)

	postgresConn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s sslrootcert=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
		config.Database.User,
		config.Database.Password,
		config.Database.SSLMode,
		config.Database.SSLRootCert,
	)

	db, err := sql.Open("postgres", postgresConn)
	if err != nil {
		panic(err)
	}

	DB = db
}

func closeDB() {
	err := DB.Close()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	defer closeDB()
	telegramClient.Poll(time.Millisecond*100, HandleUpdate)
}
