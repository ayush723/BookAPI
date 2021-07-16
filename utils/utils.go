package utils

import (
	"database/sql"
	"log"
)

func Initialize(db *sql.DB){
	const book = `
	create table if not exists books_new(
		ID int primary key,
		Name varchar(40),
		Author varchar(40),
		Year varchar(40)
	)
	`

	statement, driverErr := db.Prepare(book)
	if driverErr != nil{
		log.Println(driverErr)
	}
	//create book table
	_, statementErr := statement.Exec()
	if statementErr !=nil{
		log.Println("table already exists!")
	}
	log.Println("tables initialized successfully")

}