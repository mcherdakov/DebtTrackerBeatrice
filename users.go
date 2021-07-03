package main

import (
	"database/sql"
	"log"
)

type UserTable struct {
	id               sql.NullInt32
	username         string
	chat_id          int
	default_reciever sql.NullInt32
}

func (user *UserTable) Insert() {
	err := DB.QueryRow(
		"INSERT INTO users(username, chat_id) VALUES ($1, $2) RETURNING id",
		user.username,
		user.chat_id,
	).Scan(&user.id)
	if err != nil {
		log.Fatalln(err)
	}
}

func (user *UserTable) Update() {
	_, err := DB.Exec(
		"UPDATE users SET username=$1, chat_id=$2, default_reciever=$3 WHERE id=$4",
		user.username, user.chat_id, user.default_reciever, user.id,
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func GetOrCreateUser(telegramUser TelegramUser) (UserTable, bool) {
	user := UserTable{
		username: telegramUser.Username,
		chat_id:  telegramUser.Id,
	}

	err := DB.QueryRow(
		"SELECT id, default_reciever FROM users WHERE username=$1",
		user.username,
	).Scan(&user.id, &user.default_reciever)

	created := false
	if err == sql.ErrNoRows {
		user.Insert()
		created = true
	} else if err != nil {
		log.Fatalln(err)
	}

	return user, created
}

func GetUserByUsername(username string) (UserTable, error) {
	user := UserTable{
		username: username,
	}

	err := DB.QueryRow(
		"SELECT id, default_reciever, chat_id FROM users WHERE username=$1",
		user.username,
	).Scan(&user.id, &user.default_reciever, &user.chat_id)

	return user, err
}

func GetUserById(id sql.NullInt32) (UserTable, error) {
	user := UserTable{
		id: id,
	}

	err := DB.QueryRow(
		"SELECT default_reciever, username, chat_id FROM users WHERE id=$1",
		user.id,
	).Scan(&user.default_reciever, &user.username, &user.chat_id)

	return user, err
}
