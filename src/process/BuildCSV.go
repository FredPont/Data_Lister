// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Written by Frederic PONT.
//(c) Frederic Pont 2023

package process

import (
	conf "Data_Lister/src/configuration"
	"Data_Lister/src/pogrebdb"
	"Data_Lister/src/types"
	"fmt"
	"log"
	"strings"

	"github.com/akrylysov/pogreb"
)

func WriteCSV(fDB, dtDB *pogreb.DB, pref types.Conf) {
	header := []string{"Path", "Name", "Size", "LastAccessDate", "DirType", "TypeScore"}
	userCols, defaultValues := conf.ReadOptionalColumns()
	header = append(header, userCols...)
	fmt.Println(strings.Join(header, "\t"))
	userValues := strings.Join(defaultValues, "\t")

	it := fDB.Items()
	for {
		var dirInfo []byte
		key, val, err := it.Next()
		if err == pogreb.ErrIterationDone {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if pref.GuessDirType {
			dirInfo = pogrebdb.GetKeyDB(dtDB, key)
			if dirInfo == nil || !pref.GuessDirType {
				dirInfo = pogrebdb.StringToByte(" \t ")
			}
		}
		//log.Printf("%s %s", ByteToString(key), ByteToString(val))
		//log.Println(pogrebdb.ByteToString(key), pogrebdb.ByteToString(val))
		line := strings.Join([]string{pogrebdb.ByteToString(key), pogrebdb.ByteToString(val), pogrebdb.ByteToString(dirInfo), userValues}, "\t")
		fmt.Println(line)
	}
}
