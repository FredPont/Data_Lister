/*
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

 Written by Frederic PONT.
 (c) Frederic Pont 2023
*/

package merge

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
)

// read header of table
func readHeader(path string) []string {

	csvFile, err := os.Open(path)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	record, err := reader.Read() // read first line

	return record
}

// ReadCSVrowNames get the first column of a CSV file
func ReadCSVrowNames(path string) []string {
	var rowNames []string
	// Open the CSV file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a csv.Reader from the file
	reader := csv.NewReader(file)
	reader.Comma = '\t'

	// Loop over the records
	for {
		// Read the next record
		record, err := reader.Read()
		// Check for errors
		if err != nil {
			// Break the loop if the end of the file is reached
			if err == io.EOF {
				break
			}
			// Log any other error and continue
			log.Println(err)
			continue
		}
		// get the fist cell
		//fmt.Println(record)
		rowNames = append(rowNames, record[0])
	}

	return rowNames
}

// NewRowNames get the rownames that are new in the newfile
func NewRowNames(oldRows, newRows []string) map[string]struct{} {
	//var newRowNames []string
	// put the old rownames in a map
	oldMap := make(map[string]struct{}, len(oldRows))
	newMap := make(map[string]struct{}, 0)
	for _, item := range oldRows {
		oldMap[item] = struct{}{}
	}

	// find the rownames in newRows that are not in oldRows
	for _, item := range newRows {
		_, ok := oldMap[item]

		if ok {
			continue
		} else {
			newMap[item] = struct{}{}
		}
	}

	return newMap

}
