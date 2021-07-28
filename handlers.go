package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mcherdakov/telegoat"
)

func checkReceiverIsSet(user UserTable) bool {
	if !user.defaultReceiver.Valid {
		err := telegramClient.SendMessage(
			`Please set your receiver using "default" command:
			/default <username>`,
			user.chatId,
		)
		if err != nil {
			log.Println(err)
		}
	}

	return user.defaultReceiver.Valid
}

func HandleUpdate(update telegoat.Update) {
	log.Printf("message: %s, username: %s", update.Message.Text, update.Message.From.Username)

	user, created := GetOrCreateUser(update.Message.From)

	commandArray := strings.Fields(update.Message.Text)

	if len(commandArray) == 0 {
		return
	}

	if created {
		err := telegramClient.SendMessage(
			"Looks like you new here, hello!",
			user.chatId,
		)
		if err != nil {
			log.Println(err)
		}
	}

	if commandArray[0] != "/default" && !checkReceiverIsSet(user) {
		return
	}

	amount, err := strconv.ParseFloat(commandArray[0], 64)
	if err == nil {
		addTransaction(user, amount, strings.Join(commandArray[1:], " "))
		err := telegramClient.SendMessage("Done <3", user.chatId)
		if err != nil {
			log.Println(err)
		}
		return
	}

	switch commandArray[0] {
	case "/default":
		if len(commandArray) < 2 {
			err := telegramClient.SendMessage(
				`Correct way to do it: 
				/default <username>`,
				user.chatId,
			)
			if err != nil {
				log.Println(err)
			}
			return
		}
		setDefaultReceiver(user, commandArray[1])
	case "/d":
		getDebt(user)
	default:
		err := telegramClient.SendMessage(
			`Command not found
			Available commands:
			/default <username>
			/d
			<amount> <message>`,
			user.chatId,
		)
		if err != nil {
			log.Println(err)
		}
	}
}

func setDefaultReceiver(user UserTable, receiverUsername string) {
	receiver, err := GetUserByUsername(receiverUsername)
	if err == sql.ErrNoRows {
		err := telegramClient.SendMessage(
			"No such user in our database! Feel free to invite",
			user.chatId,
		)
		if err != nil {
			log.Println(err)
		}
		return
	} else if err != nil {
		log.Fatalln(err)
	}

	user.defaultReceiver = receiver.id
	user.Update()

	err = telegramClient.SendMessage("Done!", user.chatId)
	if err != nil {
		log.Println(err)
	}
}

func addTransaction(user UserTable, amount float64, message string) {
	userFrom := user.id.Int32
	userTo := user.defaultReceiver.Int32

	transaction := TransactionTable{
		amount:               amount,
		userFrom:             userFrom,
		userTo:               userTo,
		message:              message,
		transactionTimestamp: time.Now().Unix(),
	}
	transaction.Insert()

	dtFrom := DebtTable{
		userFrom: userFrom,
		userTo:   userTo,
	}
	dtTo := DebtTable{
		userFrom: userTo,
		userTo:   userFrom,
	}

	dtFrom.UpdateOrCreate(amount)
	dtTo.UpdateOrCreate(-amount)
}

func getDebt(user UserTable) {
	dt, err := GetDebtByUser(user)
	if err == sql.ErrNoRows {
		err := telegramClient.SendMessage(
			"No transaction yet",
			user.chatId,
		)
		if err != nil {
			log.Println(err)
		}
		return
	} else if err != nil {
		log.Fatalln(err)
	}

	receiver, err := GetUserById(user.defaultReceiver)
	if err != nil {
		log.Fatalln(err)
	}

	err = telegramClient.SendMessage(
		fmt.Sprintf(
			"Your debt to %s is %f",
			receiver.username,
			-dt.amount,
		),
		user.chatId,
	)
	if err != nil {
		log.Println(err)
	}
}
