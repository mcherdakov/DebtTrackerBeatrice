package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func SendMessage(message string, chat_id int) {
	sendMessageContext := map[string]interface{}{
		"chat_id": chat_id,
		"text":    message,
	}

	postRequest("sendMessage", sendMessageContext)
}

func checkRecieverIsSet(user UserTable) bool {
	if !user.default_reciever.Valid {
		SendMessage(
			`Please set your reciever using "default" command:
			/default <username>`,
			user.chat_id,
		)
	}

	return user.default_reciever.Valid
}

func HandleUpdate(update Update) {
	log.Printf("message: %s, username: %s", update.Message.Text, update.Message.From.Username)

	user, created := GetOrCreateUser(update.Message.From)

	commandArray := strings.Fields(update.Message.Text)

	if len(commandArray) == 0 {
		return
	}

	if created {
		SendMessage(
			"Looks like you new here, hello!",
			user.chat_id,
		)
	}

	if commandArray[0] != "/default" && !checkRecieverIsSet(user) {
		return
	}

	amount, err := strconv.ParseFloat(commandArray[0], 64)
	if err == nil {
		addTransaction(user, amount, strings.Join(commandArray[1:], " "))
		SendMessage("Done <3", user.chat_id)
		return
	}

	switch commandArray[0] {
	case "/default":
		if len(commandArray) < 2 {
			SendMessage(
				`Correct way to do it: 
				/default <username>`,
				user.chat_id,
			)
			return
		}
		setDefaultReciever(user, commandArray[1])
	case "/d":
		getDebt(user)
	default:
		SendMessage(
			`Command not found
			Available commands:
			/default <username>
			/d
			<amount> <message>`,
			user.chat_id,
		)
	}
}

func setDefaultReciever(user UserTable, recieverUsername string) {
	reciever, err := GetUserByUsername(recieverUsername)
	if err == sql.ErrNoRows {
		SendMessage(
			"No such user in our database! Feel free to invite",
			user.chat_id,
		)
		return
	} else if err != nil {
		log.Fatalln(err)
	}

	user.default_reciever = reciever.id
	user.Update()

	SendMessage("Done!", user.chat_id)
}

func addTransaction(user UserTable, amount float64, message string) {
	user_from := user.id.Int32
	user_to := user.default_reciever.Int32

	transaction := TransactionTable{
		amount:                amount,
		user_from:             user_from,
		user_to:               user_to,
		message:               message,
		transaction_timestamp: time.Now().Unix(),
	}
	transaction.Insert()

	dtFrom := DebtTable{
		user_from: user_from,
		user_to:   user_to,
	}
	dtTo := DebtTable{
		user_from: user_to,
		user_to:   user_from,
	}

	dtFrom.UpdateOrCreate(amount)
	dtTo.UpdateOrCreate(-amount)
}

func getDebt(user UserTable) {
	dt, err := GetDebtByUser(user)
	if err == sql.ErrNoRows {
		SendMessage(
			"No transaction yet",
			user.chat_id,
		)
		return
	} else if err != nil {
		log.Fatalln(err)
	}

	reciever, err := GetUserById(user.default_reciever)
	if err != nil {
		log.Fatalln(err)
	}

	SendMessage(
		fmt.Sprintf(
			"Your debt to %s is %f",
			reciever.username,
			-dt.amount,
		),
		user.chat_id,
	)
}
