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

	CreateSQLiteDB(tableName, DBpath, userCols)
}

// CreateSQLiteDB build an empty SQLite DB with primary key based on the user optional column
func CreateSQLiteDB(tableName, DBpath string, optionalColumns []string) bool {

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
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS ` + tableName + `(
	id INTEGER PRIMARY KEY,
	Path TEXT UNIQUE, Name TEXT, Modified TEXT, Size INTEGER, DirType TEXT, TypeScore REAL
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

// InsertRecord insert one row in the DataBase
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
	//fmt.Println("Connected to the database")

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

func PrepareRecord(tableName, DBpath string, fDB, dtDB, dsizeDB *pogreb.DB, pref types.Conf) {

	userColNames, defaultValues := conf.ReadOptionalColumns()
	//  =======================================================
	// read the dir/files infos stored in the pogreb databases
	//  =======================================================
	it := fDB.Items()
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
			dirInfo = pogrebdb.GetKeyDB(dtDB, key)
			if dirInfo == nil || !pref.GuessDirType {
				dirInfo = pogrebdb.StringToByte("\t")
			}
		}

		if pref.CalcSize {
			dirSize = pogrebdb.GetKeyDB(dsizeDB, key)
		}

		Path := pogrebdb.ByteToString(key)
		NameDate := strings.Split(pogrebdb.ByteToString(val), "\t")
		Name, Modified := NameDate[0], NameDate[1]
		Size := pogrebdb.ByteToInt(dirSize)
		dirTypeScore := strings.Split(pogrebdb.ByteToString(dirInfo), "\t")
		DirType, TypeScore := dirTypeScore[0], dirTypeScore[1]

		rec := []any{Path, Name, Modified, Size, DirType, TypeScore}
		// create a new slice of any with the same length as defaultValues
		strSlice := make([]any, len(defaultValues))

		// loop over strs and convert each string to an interface value
		for i := range strSlice {
			strSlice[i] = defaultValues[i]
		}

		rec = append(rec, strSlice...)

		InsertRecord(tableName, DBpath, rec, userColNames)
		//InsertRecord(tableName, DBpath, []any{key, val, dirInfo, dirSize})
	}
}
