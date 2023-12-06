package process

import (
	"fmt"
	"os"
	"path/filepath"

	getFolderSize "github.com/markthree/go-get-folder-size/src"
)

// DirSize1 return the directory size in Mbytes : quick but unstable
func DirSize1(dir string) int64 {
	size, err := getFolderSize.Invoke(dir) // get the size in bytes
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(size / 1e6) // print the size in Mbytes
	//return (size / 1e6) // return the size in Mbytes
	return (size)
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
