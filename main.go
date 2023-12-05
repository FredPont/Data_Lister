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
	"Data_Lister/src/process"
	"fmt"
	"time"
)

func main() {
	process.Title()
	t0 := time.Now()
	fmt.Println("Starting directory analysis...")
	// start a new goroutine that runs the spinner function
	// Create a channel called stop
	stop := make(chan struct{})
	go process.Spinner(stop) // enable spinner

	process.Parse()

	close(stop) // closing the channel stop the goroutine
	t1 := time.Now()
	fmt.Println("\ndone !")
	fmt.Println("Elapsed time : ", t1.Sub(t0))
}
