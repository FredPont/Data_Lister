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

package merge

import (
	"fmt"
	"reflect"
)

func CheckHeader(header1, header2 []string) bool {
	// check the headers
	if !reflect.DeepEqual(header1, header2) {
		fmt.Println("The header of new table is different from the header of the old table. The tables cannot be merged !")
		return false
	}
	return true
}
