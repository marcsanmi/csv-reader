package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

type CustomerModel struct {
	DB *sql.DB
}

func (m *CustomerModel) InsertBulkTransaction(entries chan []string) {
	stmt := `INSERT INTO customers (id, first_name, last_name, email, phone)
			VALUES (?, ?, ?, ?, ?)`

	tx, err := m.DB.Begin()

	var counter int
	for e := range entries {
		fmt.Println("e: ", e)
		if counter != 2000 {
			fmt.Println("===> executing e[1]: ", e[1])
			_, _ = tx.Exec(stmt, e[0], e[1], e[2], e[3] ,e[4])
		}

		if counter == 1000 {
			tx.Commit()
			counter = 0
			//begin another txn
			tx, err = m.DB.Begin()
		}
		counter = counter + 1
	}

	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}