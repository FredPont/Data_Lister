package process

import (
	conf "Data_Lister/src/configuration"
	"Data_Lister/src/types"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

func Parse() {
	dirSignatures := conf.ReadDirSignatures()
	fmt.Println(dirSignatures)
	pref := conf.ReadConf()
	rootLevel := strings.Count(pref.InputDir, string(os.PathSeparator))
	err := readDir(pref.InputDir, rootLevel, dirSignatures, pref)
	if err != nil {
		panic(err)
	}
}

func readDir(path string, rootLevel int, dirSignatures map[string]types.DirSignature, pref types.Conf) error {
	//fmt.Println("level=", rootLevel)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	names, _ := file.Readdirnames(0)
	fmt.Println(names)
	dirScore := ScoreType(names, dirSignatures)
	if dirScore.IsMatch {
		fmt.Println(path, " -> ", dirScore.Label, dirScore.Score)
		return nil
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
		DirOutput(filePath, "file", fileInfo, pref)
		if fileInfo.IsDir() {
			DirOutput(filePath, "dir", fileInfo, pref)

			if Level(filePath, rootLevel) < pref.Level {
				readDir(filePath, rootLevel, dirSignatures, pref)
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
			reg := regexp.QuoteMeta(signature)
			for _, name := range names {
				result, err := regexp.MatchString(reg, name)
				if err != nil {
					fmt.Println(err)
				}
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
				//fmt.Println("trouvé!", label, score, dirSig.ScoreThreshold)
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
func DirOutput(filePath, fileType string, info fs.FileInfo, pref types.Conf) {
	if pref.ListFiles && fileType == "file" {
		saveOutput(filePath, info)
	} else if fileType == "dir" {
		saveOutput(filePath, info)
	}
}

// saveOutput save the file/dir information
func saveOutput(filePath string, info fs.FileInfo) {
	fmt.Println(filePath, "size=", info.Size())
}
