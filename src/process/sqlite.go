package process

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func sqlite() {
	tableName := "id_table"

	// Open the database connection
	db, err := sql.Open("sqlite3", "./test2.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database")

	// Create a table
	// 	sqlStmt := `
	// CREATE TABLE IF NOT EXISTS users (
	// id INTEGER PRIMARY KEY,
	// name TEXT,
	// email TEXT
	// );
	// `
	// 	_, err = db.Exec(sqlStmt)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println("Created table users")

	// Insert some rows
	sqlStmt := `
INSERT INTO ` + tableName + ` (Name, DirType) VALUES (?, ?);
`
	_, err = db.Exec(sqlStmt, "FredTest1", "BCL2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a row into users")

	_, err = db.Exec(sqlStmt, "FredTest2", "Image")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted another row into users")
}
