package process

import (
	conf "Data_Lister/src/configuration"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitSQL read the DBpath from json and the user optional columns from CSV before creating a SQL database
func InitSQL() {
	pref := conf.ReadConf() // read preferences
	DBpath := pref.OutputFile
	userCols, _ := conf.ReadOptionalColumns()
	log.Println(DBpath, userCols)
	CreateSQLiteDB("data", DBpath, userCols)
}

// CreateSQLiteDB build an empty SQLite DB with primary key based on the user optional column
func CreateSQLiteDB(tableName, DBpath string, optionalColumns []string) bool {

	// Open the database connection
	db, err := sql.Open("sqlite3", DBpath)
	if err != nil {
		//log.Fatal(err)
		log.Println(DBpath, err)
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database")

	// optional cols are set as TEXT
	var optCol string
	for _, col := range optionalColumns {
		optCol = optCol + "," + col + " TEXT"
	}

	// Create a table
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS ` + tableName + `(
	id INTEGER PRIMARY KEY,
	Path TEXT, Name TEXT, Modified TEXT, Size INTEGER, DirType TEXT, TypeScore REAL
	` + optCol + `
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created table " + tableName)

	return true
}

func InsertRecord(tableName, DBpath string, records []any) bool {

	// Open the database connection
	db, err := sql.Open("sqlite3", DBpath)
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

	// Insert some rows
	sqlStmt := `
INSERT INTO ` + tableName + ` (Path, Name, Modified, Size, DirType, TypeScore) VALUES (?, ?, ?, ?, ?, ?);
`
	_, err = db.Exec(sqlStmt, records...)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Inserted a row into table")
	return true
}
