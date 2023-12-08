package process

import (
	"Data_Lister/src/pogrebdb"
	"fmt"
	"os"
	"path/filepath"

	"github.com/akrylysov/pogreb"
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

func DirSize2(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}
		fmt.Println(path, info.Name(), info.Size(), size)
		return err
	})
	fmt.Println("size returned=", size, path)
	return size, err
}

func DirSize(path string, dsizeDB *pogreb.DB) {
	rootDir := path
	// Create a map to store the subdirectory sizes
	sizes := make(map[string]int64)

	// Define the callback function that will be called for each file or folder
	walkFn := func(path string, info os.FileInfo, err error) error {
		// If there is an error, return it
		if err != nil {
			return err
		}
		// If the file is a folder, add its name to the map with zero size
		if info.IsDir() {
			sizes[path] = 0
		} else {
			// If the file is not a folder, add its size to all the parent folder sizes in the map
			//dir := filepath.Dir(path)
			//sizes[dir] += info.Size()
			incrementDirSizes(rootDir, path, info.Size(), sizes)
		}
		// Return nil to continue the walk
		return nil
	}

	// Walk through the current directory and its subdirectories with the callback function
	err := filepath.Walk(path, walkFn)
	// If there is an error, print it and exit the program
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the map of subdirectory sizes
	//fmt.Println(sizes)

	// enter the dirpath => sizes into the dsizeDB
	pogrebdb.MapToDB(sizes, dsizeDB)
}

// incrementDirSizes propagate current file size to all upper dir up to root dir
func incrementDirSizes(rootDir, path string, fileSize int64, sizes map[string]int64) {
	//rootDir := filepath.Dir(path)
	for {
		path = filepath.Dir(path)
		sizes[path] += fileSize
		//fmt.Println(path, rootDir)
		if path == rootDir {
			return
		}

	}

}
