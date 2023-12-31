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

package process

import (
	conf "Data_Lister/src/configuration"
	"Data_Lister/src/pogrebdb"
	"Data_Lister/src/types"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/akrylysov/pogreb"
)

func WriteCSV(outputFile string, fDB, dtDB, dsizeDB *pogreb.DB, pref types.Conf) {
	//fmt.Println("writing results to ", outputFile)
	//  =========================
	// build result table header
	//  =========================
	header := []string{"Path", "Name", "Modified", "Size", "DirType", "TypeScore"}
	userCols, defaultValues := conf.ReadOptionalColumns()
	header = append(header, userCols...)

	// userValues string storing the default values of user custom columns
	//userValues := strings.Join(defaultValues, "\t")

	// Create a file to write to
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Create a new csv.Writer
	writer := csv.NewWriter(file)
	// Set the delimiter to tab
	writer.Comma = '\t'
	// write result table header
	writeLine(writer, header)

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

		writeLine(writer, formatOutput(key, val, dirInfo, dirSize, defaultValues))
	}
	// Flush the buffered data
	writer.Flush()
}

func writeLine(writer *csv.Writer, data []string) {
	// Write the []string as a row to the file
	err := writer.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func formatOutput(key, val, dirInfo, dirSize []byte, defaultValues []string) []string {
	out := []string{pogrebdb.ByteToString(key)}
	out = append(out, strings.Split(pogrebdb.ByteToString(val), "\t")...)
	out = append(out, pogrebdb.IntToString(pogrebdb.ByteToInt(dirSize)))
	out = append(out, strings.Split(pogrebdb.ByteToString(dirInfo), "\t")...)

	out = append(out, defaultValues...)
	return out
}
