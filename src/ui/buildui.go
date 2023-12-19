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
	conf "Data_Lister/src/configuration"
	"Data_Lister/src/merge"
	"Data_Lister/src/types"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Regist contains the software interface
type Regist struct {
	win          fyne.Window
	img          fyne.CanvasObject
	cardTitle    string
	cardSubTitle string
	config       types.Conf
}

// NewRegist create a new registration application
func NewRegist() *Regist {

	im := picture()
	pref := conf.ReadConf() // read preferences

	regist := &Regist{img: im, config: pref}
	return regist
}

// BuildUI creates the main window of our application
func (reg *Regist) BuildUI(w fyne.Window) {
	reg.win = w

	// user settings
	var userSetting types.Conf
	//---------
	// home tab
	progBar := widget.NewProgressBarInfinite()
	progBar.Hide()

	inputDirURL := binding.NewString()
	inputDirURL.Set(insertNewlines(reg.config.InputDir, 45))
	inputDirStr, _ := inputDirURL.Get()
	//inputDirLabel := widget.NewLabel(inputDirStr)
	inputDirLabel := widget.NewLabelWithStyle(inputDirStr, fyne.TextAlignLeading, fyne.TextStyle{})
	inputDirButton := getdirPath(reg.win, "Choose the directory to scan", inputDirURL, inputDirLabel)

	//oldFileURL.Set(oldfile)
	// Create a string binding
	outFileURL := binding.NewString()
	outFileURL.Set(insertNewlines(reg.config.OutputFile, 45))
	outFileStr, _ := outFileURL.Get()
	outFileLabel := widget.NewLabelWithStyle(outFileStr, fyne.TextAlignLeading, fyne.TextStyle{})
	outFileButton := getfileSave(reg.win, "Output file", outFileURL, outFileLabel) // label a dissocier
	//outFileButton := getfilePath(reg.win, "Output file", outFileURL, outFileLabel)

	// image logo
	pict := widget.NewCard(reg.cardTitle, reg.cardSubTitle, reg.img)

	listfiles := widget.NewCheck("List Files", func(v bool) {})
	listfiles.Checked = false // set the default value to false

	guessType := widget.NewCheck("Guess Dir Type", func(v bool) {})
	guessType.Checked = true // set the default value to true

	dirSize := widget.NewCheck("Compute dir Size", func(v bool) {})
	dirSize.Checked = false // set the default value to false

	levelLab := widget.NewLabel("Level")
	level := widget.NewEntry()
	level.SetText(IntToString(reg.config.Level))
	//level.SetPlaceHolder("3")
	levelEntry := container.New(layout.NewHBoxLayout(), levelLab, level)

	//-------------
	// filters tab
	includeRegex := widget.NewCheck("Include Regex", func(v bool) {})
	includeRegex.Checked = false // set the default value to true

	var includeFormated []string
	include := widget.NewMultiLineEntry()
	include.OnChanged = func(s string) {
		includeFormated = textValidation(s)
		log.Println(includeFormated)
	}

	excludeRegex := widget.NewCheck("Exclude Regex", func(v bool) {})
	excludeRegex.Checked = false // set the default value to true
	var excludeFormated []string
	exclude := widget.NewMultiLineEntry()
	exclude.OnChanged = func(s string) {
		excludeFormated = textValidation(s)
		log.Println(includeFormated)
		log.Println(excludeFormated)
	}

	dateFilter := widget.NewCheck("Date Filter", func(v bool) {})
	dateFilter.Checked = false // set the default value to true

	olderthan := widget.NewEntry()
	olderthan.SetPlaceHolder("3023-12-12")
	olderthan.Validator = dateValidator

	newerthan := widget.NewEntry()
	newerthan.SetPlaceHolder("1922-12-12")
	newerthan.Validator = dateValidator

	filtersContent := container.NewVBox(
		includeRegex,
		include,
		excludeRegex,
		exclude,
		dateFilter,
		olderthan,
		newerthan,
	)

	//---------
	// merge

	// Create a string binding
	oldFileURL := binding.NewString()
	oldFileButton := getfilePath(reg.win, "Choose file to update", oldFileURL)
	//oldFileButton

	//oldFileURL.Set(oldfile)
	// Create a string binding
	newFileURL := binding.NewString()
	newFileButton := getfilePath(reg.win, "Choose new file", newFileURL)
	mergeButton := widget.NewButton("Merge Files", func() {
		progBar.Show()
		oldURL, _ := oldFileURL.Get()
		newURL, _ := newFileURL.Get()
		merge.Merge(cleanFileURL(oldURL), cleanFileURL(newURL))
		//time.Sleep(3 * time.Second) // pauses the execution for 3 seconds
		log.Println(oldURL, newURL)
		progBar.Hide()
	})
	mergeContent := container.NewVBox(oldFileButton, newFileButton, mergeButton, progBar)
	// Create a widget label with some help text
	helpContent := container.NewVScroll(widget.NewLabel(helpText()))

	// buttons
	closeButton := widget.NewButtonWithIcon("Close", theme.LogoutIcon(), func() { reg.win.Close() })
	runButton := widget.NewButtonWithIcon("Run", theme.ComputerIcon(), func() {
		userSetting = reg.GetUserSettings(inputDirURL, outFileURL,
			listfiles, guessType, dirSize, includeRegex, excludeRegex, dateFilter, level)
		log.Println(userSetting)
		reg.saveConfig(userSetting)
		//reg.win.Close()
	})

	//homeContent := container.NewVBox(listfiles, guessType, dirSize, levelEntry, closeButton, pict, progBar)
	homeContent := container.New(layout.NewGridLayoutWithColumns(2),
		container.NewVBox(inputDirButton, inputDirLabel, outFileButton, outFileLabel, listfiles, guessType, dirSize, levelEntry, runButton, closeButton, progBar),
		pict)

	//--------------------
	// build windows tabs
	homeTab := container.NewTabItem("Home", homeContent)
	filtersTab := container.NewTabItem("Filters", filtersContent)
	mergeTab := container.NewTabItem("Merge", mergeContent)
	helpTab := container.NewTabItem("Help", helpContent)

	// build tab container
	tabs := container.NewAppTabs(homeTab, filtersTab, mergeTab, helpTab)

	// buttonsBlock := container.NewGridWithRows(2,
	// 	layout.NewSpacer(),
	// 	buttons,
	// )

	content := tabs

	w.SetContent(content)

	w.Content().Refresh()
}
