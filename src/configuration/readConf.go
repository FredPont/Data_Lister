package conf

import (
	"encoding/json"
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
