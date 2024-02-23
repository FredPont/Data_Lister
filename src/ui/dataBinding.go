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
)

func errorBinding(err error, window fyne.Window) {
	dialog.ShowError(err, window) // show the error dialog
	window.ShowAndRun()
}

// GetUserSettings get user settings from the gui and record them in a configuration struct
func (reg *Regist) GetUserSettings(guiParam types.GuiSettings) types.Conf {
	userSetting := types.Conf{
		InputDir:          reg.dbToStr(guiParam.InputDirURL),
		OutputFile:        reg.dbToStr(guiParam.OutFileURL),
		OutputDB:          reg.dbToStr(guiParam.SqliteOutFileURL),
		ListFiles:         guiParam.Listfiles.Checked,
		GuessDirType:      guiParam.GuessType.Checked,
		CalcSize:          guiParam.DirSize.Checked,
		Level:             StrToInt(guiParam.Level.Text),
		FilterName:        guiParam.Selection.FilterName,
		FilterPath:        guiParam.Selection.FilterPath,
		FilterPathName:    guiParam.Selection.FilterPathName,
		IncludeRegex:      guiParam.IncludeRegex.Checked,
		Include:           guiParam.IncludeFormated,
		IncludeAndExclude: guiParam.IncludeAndExclude.Checked,
		ExcludeRegex:      guiParam.ExcludeRegex.Checked,
		Exclude:           guiParam.ExcludeFormated,
		DateFilter:        guiParam.DateFilter.Checked,
		UseSQLite:         reg.dbToBool(guiParam.UseSQLiteBind),
		SQLiteTable:       reg.dbToStr(guiParam.SqlTableName),
		OlderThan:         guiParam.Olderthan.Text,
		NewerThan:         guiParam.Newerthan.Text,
	}
	return userSetting
}

// dbToStr convert dataBind binding.String to string
func (reg *Regist) dbToStr(dataBind binding.String) string {
	s, err := dataBind.Get() // s is a string, err is an error
	if err != nil {
		// handle error
		errorBinding(err, reg.win)
	}
	return s
}

// dbToBool convert dataBind binding.Bool to Bool
func (reg *Regist) dbToBool(dataBind binding.Bool) bool {
	b, err := dataBind.Get() // s is a string, err is an error
	if err != nil {
		// handle error
		errorBinding(err, reg.win)
	}
	return b
}
