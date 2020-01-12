package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	//"github.com/marcsanmi/csv-reader/pkg/models/postgres"
	"io"
	"log"
	"os"
	_ "github.com/lib/pq"
	//"github.com/marcsanmi/csv-reader/pkg/models/postgres"
)

func read(filename string, entries chan []string) {
	fmt.Println("Start reading")

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		entry, err := reader.Read()
		fmt.Println("===> Reading: ", entry)
		if err == io.EOF {
			break
		}
		entries <- entry
	}
}

func processCSV(filename string) (entries chan []string) {
	entries = make(chan []string)
	go func() {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		defer close(entries)

		reader := csv.NewReader(file)

		for {
			entry, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			entries <- entry
		}
	}()
	return
}

func start() {
	for rec := range processCSV("MOCK_DATA2.csv") {
		fmt.Println("yawwww: ", rec[1])
		fmt.Printf("%T: \n", rec)
	}
}

func convertToEntries(entry string) {

}

func main() {

	fmt.Println("Startinng main...")

	// Open db connection pool
	dsn := "docker:docker@/customers?parseTime=true"
	db, err := connectDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	// We close the pool before main exits
	defer db.Close()

	fmt.Println("No way...")

	customerModel := &CustomerModel{DB: db}
	//customerModel := &postgres.CustomerModel{DB: db}
	customerModel.InsertBulkTransaction(processCSV("MOCK_DATA2.csv"))
}

func connectDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// Check that the connection to the db is successful
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func postgresConn(uri string) *sql.DB {

	db, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("paniquinggggg")
		panic(err)
	}
	return db
}

////

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
