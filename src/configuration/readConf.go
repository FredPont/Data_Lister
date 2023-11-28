package conf

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"Data_Lister/src/types"
)

// ReadConf read json conf file
func ReadConf() types.Conf {
	fname := "config/settings.json"
	var cs types.Conf
	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	bytes, err := io.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &cs)
	if err != nil {
		panic(err)
	}
	//fmt.Println(cs)
	return cs
}

// ReadDirSignatures read Directories signature json conf file
func ReadDirSignatures() map[string]types.DirSignature {
	fname := "config/DirSignatures.json"
	cs := make(map[string]types.DirSignature, 5)
	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	bytes, err := io.ReadAll(fp)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &cs)
	if err != nil {
		panic(err)
	}
	//fmt.Println(cs)
	return cs
}

func ReadOptionalColumns() ([]string, []string) {
	// Open the CSV file
	file, err := os.Open("config/columns.tsv")
	if err != nil {
		fmt.Println(err)
		return []string{}, []string{}
	}
	// Close the file when the function returns
	defer file.Close()

	// Create a new csv.Reader
	reader := csv.NewReader(file)
	// Set the delimiter to TAB
	reader.Comma = '\t'
	// Set the number of fields per record to -1, which means variable
	reader.FieldsPerRecord = -1

	// Read the first line as the header
	_, err = reader.Read()
	if err != nil {
		fmt.Println(err)
		return []string{}, []string{}
	}

	// Create an empty slice of []string
	var header []string
	var defaultValues []string

	// Loop through the remaining lines
	for {
		// Read a line
		line, err := reader.Read()
		// Check the error value
		if err != nil {
			// Break the loop when the end of the file is reached
			if err == io.EOF {
				break
			}
			// Print the error otherwise
			fmt.Println(err)
			return []string{}, []string{}
		}

		// Append the colname to header
		header = append(header, line[0])
		// Append the value to defaultValues
		defaultValues = append(defaultValues, line[1])
	}

	return header, defaultValues
}
