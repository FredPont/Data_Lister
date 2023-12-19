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

package ui

import (
	"log"
	"strconv"
)

// StrToInt converts string to int
func StrToInt(str string) int {
	i, err := strconv.Atoi(str) // i is 42, err is nil
	if err != nil {
		// handle error
		log.Println("Level illegal value ! ", str, "cannot be converted to integer !")
	}
	return i
}

// IntToString converts int to  string
func IntToString(i int) string {
	return strconv.Itoa(i)
}
