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
	"regexp"
	"strings"
	"time"
)

// FilterName apply filters to the dir/file name
// if return is false, the name is rejected
func FilterName(name string, pref types.Conf) bool {
	if len(pref.Include) > 0 && pref.Include[0] != "" {
		//fmt.Println(pref.Include)
		return IncludeFilter(name, pref)
	} else if len(pref.Exclude) > 0 && pref.Exclude[0] != "" {
		//fmt.Println(name, ExcludeFilter(name, pref))
		return ExcludeFilter(name, pref)
	}
	return true // if no filter, any name is valid
}

// ExcludeFilter apply exclusion list to the name
func ExcludeFilter(name string, pref types.Conf) bool {
	if pref.ExcludeRegex {
		excListRegex := pref.CompiledExcludeRegex
		for _, reg := range excListRegex {
			if regexFilter(name, reg) {
				return false
			}
		}
	} else {
		excList := pref.Exclude
		for _, reg := range excList {
			//fmt.Println(name, reg, stringFilter(name, reg))
			if stringFilter(name, reg) {
				return false
			}
		}
	}
	return true
}

// IncludeFilter apply inclusion list to the name
func IncludeFilter(name string, pref types.Conf) bool {
	if pref.IncludeRegex {
		incListRegex := pref.CompiledIncludeRegex
		for _, reg := range incListRegex {
			//fmt.Println(name, reg, regexFilter(name, reg))
			if regexFilter(name, reg) {
				return true
			}
		}
	} else {
		incList := pref.Include
		for _, reg := range incList {
			if stringFilter(name, reg) {
				return true
			}
		}
	}
	return false
}

// stringFilter search "reg" string in name
func stringFilter(name, reg string) bool {
	return strings.Contains(name, reg)
}

// regexFilter returns if regex "reg" match name
func regexFilter(name string, reg *regexp.Regexp) bool {
	//re := regexp.MustCompile(reg)
	return reg.MatchString(name)
}

// PreCompileAllRegex compile include/exclude regex to save compilation timea []string to []*regexp.Regexp
func PreCompileAllRegex(pref *types.Conf) {
	if pref.ExcludeRegex {
		pref.CompiledExcludeRegex = PreCompileRegex(pref.Exclude)
		//fmt.Println(pref.CompiledExcludeRegex)
	}
	if pref.IncludeRegex {
		pref.CompiledIncludeRegex = PreCompileRegex(pref.Include)
	}

}

// PreCompileRegex compile a []string to []*regexp.Regexp to save compilation time
func PreCompileRegex(stringList []string) []*regexp.Regexp {
	regList := make([]*regexp.Regexp, len(stringList))
	for i, str := range stringList {
		regList[i] = regexp.MustCompile(str)
	}
	return regList
}

// FilterDate filter the accesstime of dir/file by date
func FilterDate(accessTime time.Time, pref types.Conf) bool {
	if !pref.DateFilter {
		return true // if date filter is not set any date is valid
	}
	if pref.OlderThan != "" && pref.NewerThan != "" {
		return Between(accessTime, pref.NewerThan, pref.OlderThan)
	} else if pref.OlderThan != "" {
		return OlderThan(accessTime, pref.OlderThan)
	} else if pref.NewerThan != "" {
		return NewerThan(accessTime, pref.NewerThan)
	}
	return true // if date filter is not set any date is valid
}

// OlderThan test if accessTime is older than userValue
func OlderThan(accessTime time.Time, userValue string) bool {
	return accessTime.Before(StringToTime(userValue))
}

// NewerThan test if accessTime  is newer than userValue
func NewerThan(accessTime time.Time, userValue string) bool {
	return accessTime.After(StringToTime(userValue))
}

// Between  test if accessTime t is between time1 and time2
func Between(t time.Time, time1, time2 string) bool {
	t1 := StringToTime(time1)
	t2 := StringToTime(time2)
	// Check if t is between t1 and t2
	if t.After(t1) && t.Before(t2) {
		return true
	} else {
		return false
	}
}

// StringToTime convert "2023-01-01" to time.Time
func StringToTime(str string) time.Time {
	layout := "2006-01-02"            // layout string using the reference date
	t, err := time.Parse(layout, str) // parse the str string into a time.Time object
	if err != nil {
		fmt.Println(err)
	} else {
		return t
	}
	return time.Time{}
}
