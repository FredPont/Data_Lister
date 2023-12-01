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
	"regexp"
	"strings"
)

func FilterName(name string, pref types.Conf) bool {
	if len(pref.Include) > 0 {
		return IncludeFilter(name, pref)
	} else if len(pref.Exclude) > 0 {
		//fmt.Println(name, ExcludeFilter(name, pref))
		return ExcludeFilter(name, pref)
	}
	return true // if no filter, any name is valid
}

func ExcludeFilter(name string, pref types.Conf) bool {
	excList := pref.Exclude
	excListRegex := pref.CompiledExcludeRegex
	if pref.IncludeRegex {
		for _, reg := range excListRegex {
			if regexFilter(name, reg) {
				return false
			}
		}
	} else {
		for _, reg := range excList {
			if stringFilter(name, reg) {
				return false
			}
		}
	}
	return true
}

func IncludeFilter(name string, pref types.Conf) bool {
	incList := pref.Include
	incListRegex := pref.CompiledIncludeRegex
	if pref.IncludeRegex {
		//fmt.Println(pref)
		for _, reg := range incListRegex {
			//fmt.Println(name, reg, regexFilter(name, reg))
			if regexFilter(name, reg) {
				return true
			}
		}
	} else {
		for _, reg := range incList {
			if stringFilter(name, reg) {
				return true
			}
		}
	}
	return false
}

func stringFilter(name, reg string) bool {
	return strings.Contains(name, reg)
}

func regexFilter(name string, reg *regexp.Regexp) bool {
	//re := regexp.MustCompile(reg)
	return reg.MatchString(name)
}

func PreCompileAllRegex(pref *types.Conf) {
	if pref.ExcludeRegex {
		pref.CompiledExcludeRegex = PreCompileRegex(pref.Exclude)
		//fmt.Println(pref.CompiledExcludeRegex)
	}
	if pref.IncludeRegex {
		pref.CompiledIncludeRegex = PreCompileRegex(pref.Include)
	}

}

func PreCompileRegex(stringList []string) []*regexp.Regexp {
	regList := make([]*regexp.Regexp, len(stringList))
	for i, str := range stringList {
		regList[i] = regexp.MustCompile(str)
	}
	return regList
}
