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
	"Data_Lister/src/process"
	"Data_Lister/src/types"
	"fmt"
	"log"
	"os"
)

// BackupCSV make a copy of the CSV table
func BackupCSV(pref types.Conf) {
	source := pref.OutputFile
	path, fileName := process.GetFileAndPath(source)
	dest := fmt.Sprintf("%v"+string(os.PathSeparator)+"%v", path, process.DatePrefix(fileName)) //new file path with date time prefix to filename
	process.CopyFile(source, dest)                                                              // copy file source to dest
	log.Println("CSV backup in ", dest)
}

func UpdateCSV(pref types.Conf, filesDB types.Databases) {
	oldfile := pref.OutputFile
	BackupCSV(pref)
	path, _ := process.GetFileAndPath(oldfile)
	tempFile := fmt.Sprintf("%v"+string(os.PathSeparator)+"%v", path, "tmp.tsv")
	process.WriteCSV(tempFile, filesDB, pref)
	Merge(oldfile, tempFile)
}
