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
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2/widget"
)

//###########################################

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

// Merge append new rows from newfile to oldfile
func Merge(oldfile, newfile string) bool {
	oldRows := ReadCSVrowNames(oldfile)
	newRows := ReadCSVrowNames(newfile)
	if !CheckHeader(readHeader(oldfile), readHeader(newfile)) {
		return false
	}
	newRowNames := NewRowNames(oldRows, newRows)

	// Open the old file in append mode
	file, err := os.OpenFile(oldfile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a csv writer
	writer := csv.NewWriter(file)
	writer.Comma = '\t'

	// Open the new file in read mode
	csvFile, err := os.Open(newfile)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	// skip the header
	_, err = reader.Read()
	if err != nil {
		// Log any other error and continue
		log.Println(err)
	}

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
		_, ok := newRowNames[record[0]] // test if record is a new row
		if ok {
			// Write a new row to the file

			err = writer.Write(record)
			if err != nil {
				log.Fatal(err)
			}
		}

	}

	// Flush the writer
	writer.Flush()
	return true
}

// MergeStatusDone display merging success in the status bar
func MergeStatusDone(mergeStatus *widget.Label) {

	log.Println("Merge done !")
	// Show progress in status bar
	mergeStatus.Text = "Merge done !"
	mergeStatus.Refresh()
	time.Sleep(time.Second)
	mergeStatus.Text = "Ready"
	mergeStatus.Refresh()
}

// MergeStatusFail display merging problem in the status bar
func MergeStatusFail(mergeStatus *widget.Label) {

	// Show error in status bar
	mergeStatus.Text = "The header of new table is different from the header of the old table. The tables cannot be merged !"
	mergeStatus.Refresh()
	time.Sleep(2 * time.Second)
	mergeStatus.Text = "Ready"
	mergeStatus.Refresh()
}
