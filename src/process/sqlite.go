package process

import (
	conf "Data_Lister/src/configuration"
	"Data_Lister/src/pogrebdb"
	"Data_Lister/src/types"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/akrylysov/pogreb"
	_ "github.com/mattn/go-sqlite3"
)

// InitSQL read the DBpath from json and the user optional columns from CSV before creating a SQL database
func InitSQL() {
	pref := conf.ReadConf() // read preferences
	tableName := pref.SQLiteTable
	DBpath := pref.OutputDB
	userCols, _ := conf.ReadOptionalColumns()
	log.Println("DBpath=", DBpath)
	CreateSQLiteDB(tableName, DBpath, userCols, pref)
}

// InitSQL get the DBpath from the GUI and the user optional columns from CSV before creating a SQL database
func InitSQLGUI(DBpath string) {
	pref := conf.ReadConf() // read preferences
	tableName := pref.SQLiteTable
	userCols, _ := conf.ReadOptionalColumns()
	log.Println("DBpath=", DBpath)
	CreateSQLiteDB(tableName, DBpath, userCols, pref)
}

// CreateSQLiteDB build an empty SQLite DB with primary key and columns based on the user optional column
func CreateSQLiteDB(tableName, DBpath string, optionalColumns []string, pref types.Conf) bool {

	// Open the database connection
	db, err := sql.Open("sqlite3", DBpath)
	if err != nil {
		//log.Fatal(err)
		log.Println("erreur ici : ", DBpath, err)
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
	// Path TEXT UNIQUE avoid duplicate path in database
	// DirType TEXT, TypeScore REAL columns are disabled by default. They appear is the user select the guessDirType option
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS ` + tableName + `(
	id INTEGER PRIMARY KEY,
	Path TEXT UNIQUE, Name TEXT, Modified TEXT, Size INTEGER
	` + optCol + `
	);
	`
	if pref.GuessDirType {
		sqlStmt = `
	CREATE TABLE IF NOT EXISTS ` + tableName + `(
	id INTEGER PRIMARY KEY,
	Path TEXT UNIQUE, Name TEXT, Modified TEXT, Size INTEGER, DirType TEXT, TypeScore REAL
	` + optCol + `
	);
	`
	}

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created table " + tableName)

	return true
}

// InsertRecord insert one row in the DataBase : This function is used for unit testing only
// InsertAllRecord() is used to insert all records in the sqlite databes
func InsertRecord(tableName, DBpath string, records []any, userColNames []string) bool {

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
	//fmt.Println("Connected to the database", DBpath)

	// create a string of placeholders for the values
	placeholders := strings.Repeat("?,", len(records)-1)
	placeholders = placeholders + "?" // remove the last comma

	colnames := []string{"Path", "Name", "Modified", "Size", "DirType", "TypeScore"}
	colnames = append(colnames, userColNames...)

	// create a SQL statement to insert the values into the table
	sqlStmt := fmt.Sprintf("INSERT OR IGNORE INTO %s (%s) VALUES (%s)", tableName, strings.Join(colnames, ", "), placeholders)

	_, err = db.Exec(sqlStmt, records...)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Inserted a row into table")
	return true
}

// PrepareAllRecord build each row one by one, put all raw in a slice and insert the slice in the database
func PrepareAllRecord(tableName, DBpath string, filesDB types.Databases, pref types.Conf) {
	var allRecords [][]any
	userColNames, defaultValues := conf.ReadOptionalColumns()
	//  =======================================================
	// read the dir/files infos stored in the pogreb databases
	//  =======================================================
	it := filesDB.FileDB.Items()
	for {
		dirInfo := pogrebdb.StringToByte("\t") // dirtype and size empty by default to avoid column shift if compute dir size is enabled
		dirSize := pogrebdb.StringToByte("")
		key, val, err := it.Next()
		if err == pogreb.ErrIterationDone {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if pref.GuessDirType {
			dirInfo = pogrebdb.GetKeyDB(filesDB.DirLblDB, key)
			if dirInfo == nil || !pref.GuessDirType {
				dirInfo = pogrebdb.StringToByte("\t")
			}
		}

		if pref.CalcSize {
			dirSize = pogrebdb.GetKeyDB(filesDB.DirSizeDB, key)
		}

		Path := pogrebdb.ByteToString(key)
		NameDate := strings.Split(pogrebdb.ByteToString(val), "\t")
		Name, Modified := NameDate[0], NameDate[1]
		Size := pogrebdb.ByteToInt(dirSize)

		// columns DirType, TypeScore are disabled by default
		rec := []any{Path, Name, Modified, Size}
		if pref.GuessDirType {
			dirTypeScore := strings.Split(pogrebdb.ByteToString(dirInfo), "\t")
			DirType, TypeScore := dirTypeScore[0], dirTypeScore[1]
			rec = []any{Path, Name, Modified, Size, DirType, TypeScore}
		}

		// create a new slice of any with the same length as defaultValues
		strSlice := make([]any, len(defaultValues))

		// loop over strs and convert each string to an interface value
		for i := range strSlice {
			strSlice[i] = defaultValues[i]
		}

		rec = append(rec, strSlice...)       // create one row
		allRecords = append(allRecords, rec) // append the row to allRecords
	}

	// check if the columns number inserted by the user match the database
	if getDBcolNumber(tableName, DBpath) != len(allRecords[0]) {
		log.Println("Cannot insert record ! The column number of the database (" + fmt.Sprint((getDBcolNumber(tableName, DBpath))) +
			") does not match the number of columns selected (" + fmt.Sprint(len(allRecords[0])) + ") ! consider changing the Guess dir type option.")
		return
	}
	InsertAllRecord(tableName, DBpath, allRecords, userColNames, pref)
}

// InsertRecord insert all rows in the DataBase
func InsertAllRecord(tableName, DBpath string, allRecords [][]any, userColNames []string, pref types.Conf) bool {
	var records []any
	if len(allRecords) > 0 {
		records = allRecords[0]
	} else {
		fmt.Println("unable to insert an empty row !")
		return false
	}
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
	//fmt.Println("Connected to the database", DBpath)

	// start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}
	defer tx.Rollback() // cancel transaction in case error

	// create a string of placeholders for the values
	placeholders := strings.Repeat("?,", len(records)-1)
	placeholders = placeholders + "?" // remove the last comma

	// DirType, TypeScore  columns are disabled by default. They appear is the user select the guessDirType option
	colnames := []string{"Path", "Name", "Modified", "Size"}
	if pref.GuessDirType {
		colnames = []string{"Path", "Name", "Modified", "Size", "DirType", "TypeScore"}
	}
	colnames = append(colnames, userColNames...)

	// create a SQL statement to insert the values into the table
	sqlStmt := fmt.Sprintf("INSERT OR IGNORE INTO %s (%s) VALUES (%s)", tableName, strings.Join(colnames, ", "), placeholders)

	// Prepare INSERT command
	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	for _, record := range allRecords {
		_, err = stmt.Exec(record...)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Validate transaction
	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}
	//fmt.Println("Inserted a row into table")
	return true
}

// getDBcolNumber get the columns number of the SQLite database (except the ID column) to verify that the user
// update the database with the sames options used previously especially the guess dir type
func getDBcolNumber(tableName, DBpath string) int {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", DBpath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare a query
	stmt, err := db.Prepare("SELECT * FROM " + tableName + " LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	// Print column names
	//fmt.Println(columns)

	return len(columns) - 1 // remove the ID column witch is not inserted
}
