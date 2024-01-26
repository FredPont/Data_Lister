package process

import (
	conf "Data_Lister/src/configuration"
	"Data_Lister/src/pogrebdb"
	"Data_Lister/src/types"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"github.com/akrylysov/pogreb"
	_ "github.com/mattn/go-sqlite3"
)

// InitSQL read the DBpath from json and the user optional columns from CSV before creating a SQL database
func InitSQL(outFileURL binding.String, tableName string) {
	pref := conf.ReadConf() // read preferences
	url, _ := outFileURL.Get()
	pref.OutputFile = url
	DBpath := pref.OutputFile
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

// InsertRecord insert one row in the DataBase
func InsertRecord(tableName, DBpath string, records []any, nbColsup int) bool {

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

	// create a string of placeholders for the values
	placeholders := strings.Repeat("?,", len(records))
	placeholders = placeholders[:len(placeholders)-1] // remove the last comma

	//fmt.Println(records...)
	// userColumns := records[len(records)-nbColsup] //user columns
	// var strUserCols, strUserVals string
	// for _, col := range strUserCols {
	// 	strUserCols = strUserCols + string(col) + ", "
	// 	strUserVals = strUserVals + ",? "
	// }
	// fmt.Println(userColumns, "\n", strUserCols, "\n", strUserVals)
	// Insert some rows

	// create a SQL statement to insert the values into the table
	sqlStmt := fmt.Sprintf("INSERT INTO %s VALUES (%s)", tableName, placeholders)
	fmt.Println(sqlStmt, tableName)

	//sqlStmt := `INSERT INTO ` + tableName + ` (Path, Name, Modified, Size, DirType, TypeScore) VALUES (?, ?, ?, ?, ?, ?);`

	_, err = db.Exec(sqlStmt, records...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a row into table")
	return true
}

func PrepareRecord(tableName, DBpath string, fDB, dtDB, dsizeDB *pogreb.DB, pref types.Conf) {

	_, defaultValues := conf.ReadOptionalColumns()
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

		//line := strings.Join([]string{pogrebdb.ByteToString(key), pogrebdb.ByteToString(val), pogrebdb.ByteToString(dirInfo), userValues}, "\t")
		//fmt.Println(line)

		//writeLine(writer, formatOutput(key, val, dirInfo, dirSize, defaultValues))
		fmt.Println("PrepareRecord")
		rec := []any{key, val, dirInfo, dirSize}
		// create a new slice of any with the same length as defaultValues
		strSlice := make([]any, len(defaultValues))

		// loop over strs and convert each string to an interface value
		for i, s := range strSlice {
			strSlice[i] = s
		}

		rec = append(rec, strSlice...)
		InsertRecord(tableName, DBpath, rec, len(defaultValues))
		//InsertRecord(tableName, DBpath, []any{key, val, dirInfo, dirSize})
	}
}
