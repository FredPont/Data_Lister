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
	"Data_Lister/src/process"
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
	inputDirURL.Set(reg.config.InputDir)
	inputDirStr, _ := inputDirURL.Get()
	//inputDirLabel := widget.NewLabel(inputDirStr)
	inputDirLabel := widget.NewLabelWithStyle(insertNewlines(inputDirStr, 45), fyne.TextAlignLeading, fyne.TextStyle{})
	inputDirButton := getdirPath(reg.win, "Choose the directory to scan", inputDirURL, inputDirLabel)

	//oldFileURL.Set(oldfile)
	// Create a string binding
	outFileURL := binding.NewString()
	outFileURL.Set(reg.config.OutputFile)
	outFileStr, _ := outFileURL.Get()
	outFileLabel := widget.NewLabelWithStyle(insertNewlines(outFileStr, 45), fyne.TextAlignLeading, fyne.TextStyle{})
	outFileButton := getfileSave(reg.win, "Output file", outFileURL, outFileLabel) // label a dissocier
	//outFileButton := getfilePath(reg.win, "Output file", outFileURL, outFileLabel)

	// image logo
	pict := widget.NewCard(reg.cardTitle, reg.cardSubTitle, reg.img)

	listfiles := widget.NewCheck("List Files", func(v bool) {})
	//listfiles.Checked = false // set the default value to false
	listfiles.Checked = reg.config.ListFiles

	guessType := widget.NewCheck("Guess Dir Type", func(v bool) {})
	//guessType.Checked = true // set the default value to true
	guessType.Checked = reg.config.GuessDirType

	dirSize := widget.NewCheck("Compute dir Size", func(v bool) {})
	// dirSize.Checked = false // set the default value to false
	dirSize.Checked = reg.config.CalcSize

	levelLab := widget.NewLabel("Level")
	level := widget.NewEntry()
	level.SetText(IntToString(reg.config.Level))
	//level.SetPlaceHolder("3")
	levelEntry := container.New(layout.NewHBoxLayout(), levelLab, level)

	//-------------
	// filters tab
	includeRegex := widget.NewCheck("Include : check to use Regex instead of string", func(v bool) {})
	//includeRegex.Checked = false // set the default value to true
	includeRegex.Checked = reg.config.IncludeRegex

	var includeFormated []string
	include := widget.NewMultiLineEntry()
	include.SetText(strSliceToString(reg.config.Include))
	includeFormated = reg.config.Include
	include.OnChanged = func(s string) {
		includeFormated = textValidation(s)
		//log.Println(includeFormated)
	}

	excludeRegex := widget.NewCheck("Exclude : check to use Regex instead of string", func(v bool) {})
	//excludeRegex.Checked = false // set the default value to true
	excludeRegex.Checked = reg.config.ExcludeRegex

	var excludeFormated []string
	exclude := widget.NewMultiLineEntry()
	exclude.SetText(strSliceToString(reg.config.Exclude))
	excludeFormated = reg.config.Exclude
	exclude.OnChanged = func(s string) {
		excludeFormated = textValidation(s)
		//log.Println(excludeFormated)
	}

	dateFilter := widget.NewCheck("Date Filter", func(v bool) {})
	//dateFilter.Checked = false // set the default value to true
	dateFilter.Checked = reg.config.DateFilter

	olderthan := widget.NewEntry()
	olderthan.SetText(reg.config.OlderThan)
	//olderthan.SetPlaceHolder("3023-12-12")
	olderthan.Validator = dateValidator
	olderthanLab := widget.NewLabel("Older Than")

	newerthan := widget.NewEntry()
	newerthan.SetText(reg.config.NewerThan)
	//newerthan.SetPlaceHolder("1922-12-12")
	newerthan.Validator = dateValidator
	newerthanLab := widget.NewLabel("Newer Than")

	filtersContent := container.NewVBox(
		includeRegex,
		include,
		excludeRegex,
		exclude,
		dateFilter,
		container.NewGridWithColumns(2, olderthanLab, olderthan),
		container.NewGridWithColumns(2, newerthanLab, newerthan),
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
		//log.Println(oldURL, newURL)
		progBar.Hide()
	})
	mergeContent := container.NewVBox(oldFileButton, newFileButton, mergeButton, progBar)
	// Create a widget label with some help text
	helpContent := container.NewVScroll(widget.NewLabel(helpText()))

	// buttons
	closeButton := widget.NewButtonWithIcon("Close", theme.LogoutIcon(), func() { reg.win.Close() })
	runButton := widget.NewButtonWithIcon("Run", theme.ComputerIcon(), func() {
		progBar.Show()
		userSetting = reg.GetUserSettings(inputDirURL, outFileURL,
			listfiles, guessType, dirSize, includeRegex, excludeRegex, dateFilter, level,
			olderthan, newerthan,
			includeFormated, excludeFormated)
		log.Println(userSetting)
		reg.saveConfig(userSetting)
		process.Parse()
		progBar.Hide()

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

	// build tab containerprocess.Parse()
	tabs := container.NewAppTabs(homeTab, filtersTab, mergeTab, helpTab)

	content := tabs

	w.SetContent(content)

	w.Content().Refresh()
}

func startDirAnalysis(reg *Regist, progBar *widget.ProgressBar, inputDirURL, outFileURL binding.String,
	listfiles, guessType, dirSize, includeRegex, excludeRegex, dateFilter *widget.Check,
	level, olderthan, newerthan *widget.Entry,
	includeFormated, excludeFormated []string) {

	progBar.Show()
	userSetting := reg.GetUserSettings(inputDirURL, outFileURL,
		listfiles, guessType, dirSize, includeRegex, excludeRegex, dateFilter, level,
		olderthan, newerthan,
		includeFormated, excludeFormated)
	log.Println(userSetting)
	reg.saveConfig(userSetting)
	process.Parse()
	progBar.Hide()
}
