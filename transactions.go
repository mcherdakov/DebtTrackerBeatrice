package main

import (
	"log"
	"time"
)

type TransactionTable struct {
	id                   int
	amount               float64
	userFrom             int32
	userTo               int32
	message              string
	transactionTimestamp int64
}

type DebtTable struct {
	id       int
	userFrom int32
	userTo   int32
	amount   float64
}

func (t *TransactionTable) Insert() {
	err := DB.QueryRow(
		`INSERT INTO transactions(amount, user_from, user_to, message, transaction_timestamp)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		t.amount, t.userFrom, t.userTo, t.message,
		time.Now().Unix(),
	).Scan(&t.id)

	if err != nil {
		log.Fatalln(err)
	}
}

func (dt *DebtTable) UpdateOrCreate(amount float64) {
	res, err := DB.Exec(
		"UPDATE debt SET amount=amount+$1 WHERE user_from=$2 AND user_to=$3",
		amount, dt.userFrom, dt.userTo,
	)
	if err != nil {
		log.Fatalln(err)
	}

	nRows, err := res.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	if nRows == 0 {
		dt.amount = amount
		err = DB.QueryRow(
			"INSERT INTO debt(user_from, user_to, amount) VALUES ($1, $2, $3) RETURNING id",
			dt.userFrom, dt.userTo, dt.amount,
		).Scan(&dt.id)

		if err != nil {
			log.Fatalln(err)
		}
	}
}

func GetDebtByUser(user UserTable) (DebtTable, error) {
	dt := DebtTable{
		userFrom: user.id.Int32,
		userTo:   user.defaultReceiver.Int32,
	}

	err := DB.QueryRow(
		"SELECT amount FROM debt WHERE user_from=$1 AND user_to=$2",
		dt.userFrom, dt.userTo,
	).Scan(&dt.amount)

	return dt, err
}
