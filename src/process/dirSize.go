package process

import (
	"fmt"

	getFolderSize "github.com/markthree/go-get-folder-size/src"
)

// DirSize return the directory size in Mbytes
func DirSize(dir string) int64 {
	size, err := getFolderSize.Invoke(dir) // get the size in bytes
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(size / 1e6) // print the size in Mbytes
	return (size / 1e6) // return the size in Mbytes
}
