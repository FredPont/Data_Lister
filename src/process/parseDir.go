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

package process

import (
	conf "Data_Lister/src/configuration"
	"Data_Lister/src/pogrebdb"
	"Data_Lister/src/types"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"
)

func Parse() {
	pogrebdb.InitDB()
	filesDB := pogrebdb.LoadAllDB()

	dirSignatures := conf.ReadDirSignatures() // load dir signatures

	pref := conf.ReadConf() // read preferences

	// precompilation of include/exclude regex to speed filters
	PreCompileAllRegex(&pref)

	rootLevel := strings.Count(pref.InputDir, string(os.PathSeparator))

	err := readDir(pref.InputDir, rootLevel, dirSignatures, pref, filesDB)
	if err != nil {
		//panic(err)
		fmt.Println("Error ! ", err)
		return
	}
	CSV_SQL_output(pref, filesDB)
	pogrebdb.CloseAllDB(filesDB)
}

// CSV_SQL_output start the results saving in CSV or SQLite
func CSV_SQL_output(pref types.Conf, filesDB types.Databases) {
	fmt.Println("pref.UseSQLite =", pref.UseSQLite)
	if pref.UseSQLite {
		log.Println("start SQLite output")
		PrepareAllRecord(pref.SQLiteTable, pref.OutputDB, filesDB, pref)
		log.Println("SQLite saved in", pref.SQLiteTable, pref.OutputDB)
	} else {
		log.Println("start CSV output")
		WriteCSV(pref.OutputFile, filesDB, pref)
		log.Println("CSV saved in", pref.OutputFile)
	}
}

// readDir recursive function to read dir and files to a certain level of deepness
func readDir(path string, rootLevel int, dirSignatures map[string]types.DirSignature, pref types.Conf, filesDB types.Databases) error {
	//fmt.Println("level=", rootLevel)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	names, _ := file.Readdirnames(0)
	//fmt.Println(names)
	if pref.GuessDirType {
		dirScore := ScoreType(names, dirSignatures)
		if dirScore.IsMatch {
			//fmt.Println(path, " -> ", dirScore.Label, dirScore.Score)
			outString := strings.Join([]string{dirScore.Label, strconv.FormatFloat(dirScore.Score, 'f', -1, 64)}, "\t")
			pogrebdb.InsertDataDB(filesDB.DirLblDB, pogrebdb.StringToByte(path), pogrebdb.StringToByte(outString))
			return nil
		}
	}

	for _, name := range names {
		// if !FilterName(path, name, pref) {
		// 	continue
		// }
		filePath := fmt.Sprintf("%v"+string(os.PathSeparator)+"%v", path, name)
		//fmt.Println(filePath)
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		// store file info
		if pref.ListFiles {
			// filter the name by date
			if !fileInfo.IsDir() {
				if !FilterDate(fileInfo.ModTime(), pref) {
					continue
				} else if !FilterName(path, name, pref) {
					// do not save the file if it does not contain the include string pattern
					continue
				}
				saveOutput(filePath, fileInfo, pref, filesDB)
			}

		}
		// store dir info
		if fileInfo.IsDir() {
			//log.Println("=>", path, name)
			if FilterDate(fileInfo.ModTime(), pref) {
				if FilterName(path, name, pref) {
					// do not block subdir analysis because subdir can contain the filter string
					//log.Println(path, name)
					saveOutput(filePath, fileInfo, pref, filesDB)
					//continue
				}
			}

			// analyse subdir if current level < user level limit
			if Level(filePath, rootLevel) < pref.Level {
				readDir(filePath, rootLevel, dirSignatures, pref, filesDB)
			}
		}

	}
	return nil
}

// Level return the path deepness compared to the rootLevel, pref types.Conf
func Level(path string, rootLevel int) int {
	return strings.Count(path, string(os.PathSeparator)) - rootLevel
}

// compute a score to guess the type of content of the directory
func ScoreType(names []string, dirSignatures map[string]types.DirSignature) types.DirMatch {
	//signature := []string{".go", ".git", ".DLL", ".dll", ".r", ".jl", ".pl", "\\.json", "j[a-z]{2}n", ".+[a-z]{2}n", ".json"}
	matchNB := 0
	score := 0.
	for label, dirSig := range dirSignatures {
		for _, name := range names {
			for _, signature := range dirSig.Content {
				result := strings.Contains(name, signature)
				if result {
					matchNB++
					//fmt.Println("match!", label, matchNB, "/", len(names), signature, name)
					break
				}
			}
		}
		if len(names) > 0 { // security to avoid division by zero
			// the score is the ratio of names matching regex / number of elements in the directory
			score = float64(matchNB) / float64(len(names))
			if score >= dirSig.ScoreThreshold {

				return types.DirMatch{IsMatch: true, Label: label, Score: score}
			}
		}

	}
	return types.DirMatch{IsMatch: false, Label: "", Score: 0.}
}

// ScanDirType try to guess the type of content (=names) of the directory
func ScanDirType(names []string, pref types.Conf, dirSignatures map[string]types.DirSignature) string {
	if pref.GuessDirType {
		ScoreType(names, dirSignatures)
	}
	return ""
}

// saveOutput save the file/dir information to the pogreb databases
func saveOutput(filePath string, info fs.FileInfo, pref types.Conf, filesDB types.Databases) {
	if pref.CalcSize {
		if info.IsDir() {
			DirSize(filePath, filesDB.DirSizeDB)
		} else {
			//size = info.Size()
			pogrebdb.InsertDataDB(filesDB.DirSizeDB, pogrebdb.StringToByte(filePath), pogrebdb.IntToBytes(info.Size()))
		}
	}
	modTime := info.ModTime()
	year := strconv.Itoa(modTime.Year())
	month := fmt.Sprintf("%02d", int(modTime.Month()))
	day := fmt.Sprintf("%02d", modTime.Day())

	outString := strings.Join([]string{info.Name(), year + "-" + month + "-" + day}, "\t") // save name and date to database
	pogrebdb.InsertDataDB(filesDB.FileDB, pogrebdb.StringToByte(filePath), pogrebdb.StringToByte(outString))
}
