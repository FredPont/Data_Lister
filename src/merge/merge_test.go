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
	"testing"
)

func TestCheckHeader(t *testing.T) {
	tests := []struct {
		header1, header2 []string
		want             bool
	}{
		{[]string{"Path", "Name", "Modified", "Size", "DirType", "TypeScore"}, []string{"Path", "Name", "Modified", "Size", "DirType", "TypeScore"}, true},
		{[]string{"Path", "Name", "Modified", "Size", "DirType", "TypeScore"}, []string{"Path", "Name", "Size", "Modified", "DirType", "TypeScore"}, false},
		{[]string{"Path", "Name", "Modified", "Size", "DirType", "TypeScore"}, []string{"Path", "Name", "Modified", "DirType", "TypeScore"}, false},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := CheckHeader(tc.header1, tc.header2)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
