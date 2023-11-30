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
	"os"
	"strconv"
	"strings"

	"github.com/akrylysov/pogreb"
)

func Parse() {
	pogrebdb.InitDB()
	fDB := pogrebdb.OpenDB("db/files")
	dtDB := pogrebdb.OpenDB("db/dirTypes")
	dirSignatures := conf.ReadDirSignatures() // load dir signatures
	//fmt.Println(dirSignatures)
	pref := conf.ReadConf() // read preferences
	rootLevel := strings.Count(pref.InputDir, string(os.PathSeparator))
	err := readDir(pref.InputDir, rootLevel, dirSignatures, pref, fDB, dtDB)
	if err != nil {
		panic(err)
	}
	PreCompileAllRegex(pref)
	fmt.Println(pref.CompiledExcludeRegex)
	//pogrebdb.ShowDB(fDB)
	//pogrebdb.ShowDB(dtDB)
	WriteCSV(pref.OutputFile, fDB, dtDB, pref)
	fDB.Close()
	dtDB.Close()
}

// readDir recursive function to read dir and files to a certain level of deepness
func readDir(path string, rootLevel int, dirSignatures map[string]types.DirSignature, pref types.Conf, fDB, dtDB *pogreb.DB) error {
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
			pogrebdb.InsertDataDB(dtDB, pogrebdb.StringToByte(path), pogrebdb.StringToByte(outString))
			return nil
		}
	}

	for _, name := range names {
		filePath := fmt.Sprintf("%v/%v", path, name)
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}
		DirOutput(filePath, "file", fileInfo, pref, fDB)
		if fileInfo.IsDir() {
			DirOutput(filePath, "dir", fileInfo, pref, fDB)

			if Level(filePath, rootLevel) < pref.Level {
				readDir(filePath, rootLevel, dirSignatures, pref, fDB, dtDB)
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
		for _, signature := range dirSig.Content {
			//reg := regexp.QuoteMeta(signature)
			for _, name := range names {
				// result, err := regexp.MatchString(reg, name)
				// if err != nil {
				// 	fmt.Println(err)
				// }
				result := strings.Contains(name, signature)
				if result {
					matchNB++
					//fmt.Println("match!", label, matchNB, "/", len(names), signature, name, "sign=", i, "name=", j)
				}
			}
		}
		if len(names) > 0 {
			// the score is the ratio of names matching regex / number of elements in the directory
			score = float64(matchNB) / float64(len(names))
			if score >= dirSig.ScoreThreshold {
				if score > 1. {
					score = 1.
				}
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

// DirOutput decide the output of the file/dir information
func DirOutput(filePath, fileType string, info fs.FileInfo, pref types.Conf, fDB *pogreb.DB) {
	if pref.ListFiles && fileType == "file" {
		saveOutput(filePath, info, pref, fDB)
	} else if fileType == "dir" {
		saveOutput(filePath, info, pref, fDB)
	}
}

// saveOutput save the file/dir information to the pogreb databases
func saveOutput(filePath string, info fs.FileInfo, pref types.Conf, fDB *pogreb.DB) {
	var size int64
	if pref.CalcSize {
		if info.IsDir() {
			size = DirSize(filePath)
		} else {
			size = info.Size()
		}
	}
	modTime := info.ModTime()
	year := strconv.Itoa(modTime.Year())
	month := strconv.Itoa(int(modTime.Month()))
	day := strconv.Itoa(modTime.Day())

	//fmt.Println(filePath, "name=", info.Name(), "size=", size, "date=", year, month, day)
	outString := strings.Join([]string{info.Name(), strconv.FormatInt(size/1e3, 10), year + "-" + month + "-" + day}, "\t")
	pogrebdb.InsertDataDB(fDB, pogrebdb.StringToByte(filePath), pogrebdb.StringToByte(outString))
}
