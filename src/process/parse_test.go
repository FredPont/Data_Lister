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
	"regexp"
	"testing"
	"time"
)

func TestScoreType(t *testing.T) {
	tests := []struct {
		names         []string
		dirSignatures map[string]types.DirSignature
		want          types.DirMatch
	}{
		{[]string{"prog.pl", "data", "results", "src"}, map[string]types.DirSignature{"soft": types.DirSignature{Content: []string{".go", ".git", ".r", ".jl", ".pl"}, ScoreThreshold: 0.2}}, types.DirMatch{true, "soft", 0.25}},
		{[]string{"prog.pl", "data", "results", "src"}, map[string]types.DirSignature{"soft": types.DirSignature{Content: []string{".go", ".git", ".DLL", ".dll", ".r", ".jl", ".pl", ".json"}, ScoreThreshold: 0.2}}, types.DirMatch{true, "soft", 0.25}},
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

func TestRegexFilter(t *testing.T) {
	tests := []struct {
		names, reg string
		want       bool
	}{
		{"test.pl", ".pl", true},
		{"test.pl", "\\.pl", true},
		{"file.json", "j[a-z]{2}n", true},
		{"dir1", "dir", true},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := regexFilter(tc.names, regexp.MustCompile(tc.reg))
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}

func TestOlderThan(t *testing.T) {
	tests := []struct {
		userDate string
		modTime  time.Time
		want     bool
	}{
		{"2023-01-01", StringToTime("2023-02-02"), false},
		{"2023-01-01", StringToTime("2022-02-02"), true},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := OlderThan(tc.modTime, tc.userDate)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}

func TestNewerThan(t *testing.T) {
	tests := []struct {
		userDate string
		modTime  time.Time
		want     bool
	}{
		{"2023-01-01", StringToTime("2023-02-02"), true},
		{"2023-01-01", StringToTime("2022-02-02"), false},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := NewerThan(tc.modTime, tc.userDate)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}

func TestBetween(t *testing.T) {
	tests := []struct {
		time1, time2 string
		modTime      time.Time
		want         bool
	}{
		{"2022-01-01", "2023-02-02", StringToTime("2023-01-01"), true},
		{"2022-01-05", "2023-02-02", StringToTime("2022-01-01"), false},
		{"2022-12-12", "2023-12-12", StringToTime("2023-12-04"), true},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := Between(tc.modTime, tc.time1, tc.time2)
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}

func TestCreateSQLiteDB(t *testing.T) {
	tests := []struct {
		tableName, DBpath string
		optionalColumns   []string
		pref              types.Conf
	}{
		{"MyTable", "../../test/SQLiteTest.db", []string{"col1", "col2"}, types.Conf{}},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := CreateSQLiteDB(tc.tableName, tc.DBpath, tc.optionalColumns, tc.pref)
			if !got {
				t.Fatalf("got %v", got)
			} else {
				t.Logf("Success !")
			}

		})
	}
}

func TestInsertRecord(t *testing.T) {
	tests := []struct {
		tableName, DBpath string
		record            []any
		UserSQLcolnames   []string
		//nbColsup          int
	}{
		{"MyTable", "../../test/SQLiteTest.db", []any{"Path", "Name", "2023-01-25", 12, "fasta", 0.75, "col1", "col2"}, []string{"col1", "col2"}},
		{"MyTable", "../../test/SQLiteTest.db", []any{"Path", "Name", 2023 - 01 - 25, 12, "fasta", 0.75, "col1", "col2"}, []string{"col1", "col2"}}, // wrong date use text instead
		{"MyTable", "../../test/SQLiteTest.db", []any{"Path", "Name", "2023-01-25", 12, "bcl2", 0.8, "cells", "project1"}, []string{"col1", "col2"}},
		{"MyTable", "../../test/SQLiteTest.db", []any{"Path", "Name", "2023-01-25", 12, "bcl2", 0.8, "cells", "project1"}, []string{"col1", "col2"}},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("Index=%d", i), func(t *testing.T) {
			got := InsertRecord(tc.tableName, tc.DBpath, tc.record, tc.UserSQLcolnames)
			if !got {
				t.Fatalf("got %v", got)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
