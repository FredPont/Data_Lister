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
package process

import (
	"Data_Lister/src/types"
	"fmt"
	"reflect"
	"testing"
)

func TestScoreType(t *testing.T) {
	tests := []struct {
		names         []string
		dirSignatures map[string]types.DirSignature
		want          types.DirMatch
	}{
		{[]string{"prog.pl", "data", "results", "src"}, map[string]types.DirSignature{"soft": types.DirSignature{Content: []string{".go", ".git", ".r", ".jl", ".pl"}, ScoreThreshold: 0.2}}, types.DirMatch{true, "soft", 0.25}},
		{[]string{"prog.pl", "data", "results", "src"}, map[string]types.DirSignature{"soft": types.DirSignature{Content: []string{".go", ".git", ".DLL", ".dll", ".r", ".jl", ".pl", "\\.json", "j[a-z]{2}n", ".+[a-z]{2}n", ".json"}, ScoreThreshold: 0.2}}, types.DirMatch{true, "soft", 0.25}},
		{[]string{"sample1.fasta", "sample2.fasta.gz"}, map[string]types.DirSignature{"Fasta": types.DirSignature{Content: []string{".fasta", ".FASTA"}, ScoreThreshold: 0.8}}, types.DirMatch{true, "Fasta", 1.}},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := ScoreType(tc.names, tc.dirSignatures)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
