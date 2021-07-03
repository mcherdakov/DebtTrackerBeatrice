package main

import (
	"encoding/json"
	"log"
	"time"
)

type GetUpdatesResponse struct {
	Result []Update
}

type Update struct {
	Update_id int
	Message   Message
}

type Message struct {
	Message_id int
	From       TelegramUser
	Text       string
}

type TelegramUser struct {
	Id       int
	Username string
}

func TelegramPoll() {
	offset := 0
	for {
		time.Sleep(time.Millisecond * 100)
		updates := getUpdates(offset)

		if len(updates) == 0 {
			continue
		}

		// new offset must be greater then last update_id by 1
		offset = updates[len(updates)-1].Update_id + 1
		for _, update := range updates {
			go HandleUpdate(update)
		}
	}
}

func getUpdates(offset int) []Update {
	updateContext := map[string]interface{}{
		"offset": offset,
	}
	responseData := postRequest("getUpdates", updateContext)

	var updatesResponse GetUpdatesResponse
	error := json.Unmarshal([]byte(responseData), &updatesResponse)
	if error != nil {
		log.Fatalln(error)
	}

	return updatesResponse.Result
}
