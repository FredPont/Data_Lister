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

package types

import "regexp"

type Conf struct {
	InputDir   string `json:"InputDir"`
	OutputFile string `json:"OutputFile"`
	ListFiles  bool   `json:"ListFiles"`
	//ListDir              bool     `json:"ListDir"`
	GuessDirType         bool     `json:"GuessDirType"`
	CalcSize             bool     `json:"CalcSize"`
	Level                int      `json:"Level"`
	FilterName           bool     `json:"filterName"`
	FilterPath           bool     `json:"filterPath"`
	FilterPathName       bool     `json:"filterPathName"`
	IncludeRegex         bool     `json:"IncludeRegex"`
	Include              []string `json:"Include"`
	IncludeAndExclude    bool     `json:"IncludeAndExclude"`
	ExcludeRegex         bool     `json:"ExcludeRegex"`
	Exclude              []string `json:"Exclude"`
	OlderThan            string   `json:"OlderThan"`
	NewerThan            string   `json:"NewerThan"`
	DateFilter           bool     `json:"DateFilter"`
	UseSQLite            bool     `json:"UseSQLite"`
	SQLiteTable          string   `json:"SQLiteDB"`
	CompiledIncludeRegex []*regexp.Regexp
	CompiledExcludeRegex []*regexp.Regexp
}

type DirSignature struct {
	Content        []string `json:"content"`
	ScoreThreshold float64  `json:"scoreThreshold"`
}

type DirMatch struct {
	IsMatch bool
	Label   string
	Score   float64
}

type RadioGroupFilters struct {
	FilterName     bool
	FilterPath     bool
	FilterPathName bool
}
