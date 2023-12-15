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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func errorBinding(err error, window fyne.Window) {
	dialog.ShowError(err, window) // show the error dialog
	window.ShowAndRun()
}

// GetUserSettings get user settings from the gui
func (reg *Regist) GetUserSettings(inputDirURL, outFileURL binding.String, listFilesCB *widget.Check) types.Conf {
	userSetting := types.Conf{
		InputDir:   reg.dbToStr(inputDirURL),
		OutputFile: reg.dbToStr(outFileURL),
		ListFiles:  listFilesCB.Checked,
	}
	return userSetting
}

func (reg *Regist) dbToStr(dataBind binding.String) string {
	s, err := dataBind.Get() // s is a string, err is an error
	if err != nil {
		// handle error
		errorBinding(err, reg.win)
	}
	return s
}
