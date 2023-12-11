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

package main

import (
	"Data_Lister/src/merge"
	"Data_Lister/src/process"
	"flag"
	"fmt"
	"time"
)

func main() {
	process.Title()
	t0 := time.Now()
	cmdLine()
	//fmt.Println("Starting directory analysis...")
	// start a new goroutine that runs the spinner function
	// Create a channel called stop
	// stop := make(chan struct{})
	// go process.Spinner(stop) // enable spinner

	// process.Parse()

	// close(stop) // closing the channel stop the goroutine
	t1 := time.Now()
	fmt.Println("\ndone !")
	fmt.Println("Elapsed time : ", t1.Sub(t0))
}

func cmdLine() {
	// start DataLister in cmd line
	var cmd bool
	var mergeFiles bool
	var oldFile string //old result file to be merged with
	var newFile string // new result file
	flag.BoolVar(&mergeFiles, "m", false, "Start DataLister merging tool.")
	flag.BoolVar(&cmd, "c", false, "Start DataLister directories analysis in command line. Example : DataLister -c")
	flag.StringVar(&oldFile, "o", "", "Old result file path. Example, to add new data from newfile to oldfile : DataLister -m -o oldfile.csv -i newfile.csv")
	flag.StringVar(&newFile, "i", "", "New result file path. Only new files/dir are added to the old file")
	flag.Parse() // parse the flags

	if cmd {
		fmt.Println("Starting directory analysis...")
		// start a new goroutine that runs the spinner function
		// Create a channel called stop
		stop := make(chan struct{})
		go process.Spinner(stop) // enable spinner

		process.Parse()

		close(stop) // closing the channel stop the goroutine
	}

	if mergeFiles {
		fmt.Println("Starting update of ...", oldFile, "with", newFile)
		// start the merging tool. The old result file is merged with a new result file
		// Create a channel called stop
		stop := make(chan struct{})
		merge.Merge(oldFile, newFile)
		close(stop) // closing the channel stop the goroutine
	}

}
