package main

import (
	"log"
	"time"
)

type TransactionTable struct {
	id                    int
	amount                float64
	user_from             int32
	user_to               int32
	message               string
	transaction_timestamp int64
}

type DebtTable struct {
	id        int
	user_from int32
	user_to   int32
	amount    float64
}

func (t *TransactionTable) Insert() {
	err := DB.QueryRow(
		`INSERT INTO transactions(amount, user_from, user_to, message, transaction_timestamp)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		t.amount, t.user_from, t.user_to, t.message,
		time.Now().Unix(),
	).Scan(&t.id)

	if err != nil {
		log.Fatalln(err)
	}
}

func (dt *DebtTable) UpdateOrCreate(amount float64) {
	res, err := DB.Exec(
		"UPDATE debt SET amount=amount+$1 WHERE user_from=$2 AND user_to=$3",
		amount, dt.user_from, dt.user_to,
	)
	if err != nil {
		log.Fatalln(err)
	}

	n_rows, err := res.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	if n_rows == 0 {
		dt.amount = amount
		err = DB.QueryRow(
			"INSERT INTO debt(user_from, user_to, amount) VALUES ($1, $2, $3) RETURNING id",
			dt.user_from, dt.user_to, dt.amount,
		).Scan(&dt.id)

		if err != nil {
			log.Fatalln(err)
		}
	} else if err != nil {
		log.Fatalln(err)
	}
}

func GetDebtByUser(user UserTable) (DebtTable, error) {
	dt := DebtTable{
		user_from: user.id.Int32,
		user_to:   user.default_reciever.Int32,
	}

	err := DB.QueryRow(
		"SELECT amount FROM debt WHERE user_from=$1 AND user_to=$2",
		dt.user_from, dt.user_to,
	).Scan(&dt.amount)

	return dt, err
}
