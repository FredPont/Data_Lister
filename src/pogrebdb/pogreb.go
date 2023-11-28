// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Written by Frederic PONT.
//(c) Frederic Pont 2023

package pogrebdb

import (
	"fmt"
	"log"
	"os"

	"github.com/akrylysov/pogreb"
)

// InsertDataDB inserts one record in the database.
// Pogreb is thread safe according to https://github.com/akrylysov/pogreb#key-characteristics
func InsertDataDB(db *pogreb.DB, key, record []byte) {
	// mutex.Lock()
	// defer mutex.Unlock()
	err := db.Put([]byte(key), record)
	if err != nil {
		log.Fatal(err)
	}

}

func OpenDB(dbName string) *pogreb.DB {
	db, err := pogreb.Open(dbName, nil)
	if err != nil {
		log.Println("Cannot open the database", dbName)
		log.Fatal(err)
		//return
	}
	return db
	//defer db.Close()
}

// GetDataDB retrieve the inserted value matching the string "key" in []byte
func GetDataDB(db *pogreb.DB, key string) []byte {
	val, err := db.Get([]byte(key))
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("%s", val)
	return val
}

// GetKeyDB retrieve the inserted value matching the "key" in []byte
func GetKeyDB(db *pogreb.DB, key []byte) []byte {
	val, err := db.Get(key)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("%s", val)
	return val
}

// InitDB clear db dir and create new databases for files and dir signatures
func InitDB() {
	CleanDB("db")
	//CreateDataBase("files")
	//CreateDataBase("dirTypes")
}

// CreateDataBase create a new pogreb database in tmp dir
func CreateDataBase(dbName string) *pogreb.DB {
	db, err := pogreb.Open("db/"+dbName, nil)
	if err != nil {
		log.Fatal(err)
		fmt.Println("cannot create database in db dir !")
	}
	//defer db.Close() // do not close the database before the end of the run

	fmt.Println("database created in db/" + dbName)
	return db
}

// CleanDB remove and create a directory
func CleanDB(dir string) {
	ClearDB(dir)
	MkDir(dir)
}

// remove DB dir
func ClearDB(dir string) {
	// clear db directory
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatal(err)
	}
}

// create DB dir
func MkDir(dir string) {
	//Create a folder/directory at a full qualified path
	err := os.Mkdir(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

// ShowDB print the database content
func ShowDB(db *pogreb.DB) {
	it := db.Items()
	for {
		key, val, err := it.Next()
		if err == pogreb.ErrIterationDone {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("%s %s", ByteToString(key), ByteToString(val))
		log.Println(ByteToString(key), ByteToString(val))
	}
}
