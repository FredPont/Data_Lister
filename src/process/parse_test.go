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
