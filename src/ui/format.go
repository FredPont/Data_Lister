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
	"Data_Lister/src/types"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// textValidation convert a string containing newlines into a []string
func textValidation(cmt string) []string {
	//cmt = strings.ReplaceAll(cmt, "\n", " ")  // remove all newline
	cmt = strings.ReplaceAll(cmt, "\r", " ") // remove all carriage return
	return strings.Split(cmt, "\n")

}

// strSliceToString convert a []string to string inserting new lines.
// Used for include/exclude multiline entry of the GUI
func strSliceToString(sl []string) string {
	return strings.Join(sl, "\n")
}

// cleanFileURL removes "file://" at the beginning of URL
func cleanFileURI(url fyne.URI) string {
	//return strings.TrimPrefix(url, "file://") //
	if os.PathSeparator == '\\' {
		//fmt.Println("Dectected OS Windows")
		//return strings.Replace(url.Path(), "/", "\\", -1)
		return url.Path() // files can be accessed on window with "/" separators
	} else {
		//fmt.Println("Dectected OS Linux/Mac")
		return url.Path()
	}
}

// cleanFileURL removes "file://" at the beginning of URI and replace the separator for windows
func cleanDirURI(url fyne.ListableURI) string {
	//return strings.TrimPrefix(url, "file://") //
	if os.PathSeparator == '\\' {
		fmt.Println("Dectected OS Windows")
		return strings.Replace(url.Path(), "/", "\\", -1)
	} else {
		//fmt.Println("Dectected OS Linux/Mac")
		return url.Path()
	}
}

// dateValidator control the date format
func dateValidator(text string) error {
	const dateFormat = "2006-01-02"
	_, err := time.Parse(dateFormat, text)
	return err
}

// saveConfig export the user setting to the config/settingsjson file
func (reg *Regist) saveConfig(userSetting types.Conf) {
	fname := "config/settings.json"
	b, err := json.Marshal(userSetting) // convert the struct to JSON
	if err != nil {
		// handle error
		dialog.ShowInformation("Alert", "Cannot convert user data to struct !", reg.win) // show the alert dialog

	}
	err = os.WriteFile(fname, b, 0644) // write the JSON to a file
	if err != nil {
		// handle error
		dialog.ShowInformation("Alert", "Cannot save user data to "+fname, reg.win) // show the alert dialog
	}

}

// insertNewlines insert newline every x char to split the url in more than one line
func insertNewlines(s string, n int) string {
	var result string
	for i, runeValue := range s {
		if i > 0 && i%n == 0 {
			result += "\n"
		}
		result += string(runeValue)
	}
	return result
}
