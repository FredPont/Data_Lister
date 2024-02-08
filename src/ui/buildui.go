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
	"fmt"
	"log"
	"time"

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
	progBar      *widget.ProgressBarInfinite
	cardTitle    string
	cardSubTitle string
	config       types.Conf
}

// NewRegist create a new registration application
func NewRegist() *Regist {
	progBar := widget.NewProgressBarInfinite()
	im := picture()
	pref := conf.ReadConf() // read preferences

	regist := &Regist{img: im, progBar: progBar, config: pref}
	return regist
}

// BuildUI creates the main window of our application
func (reg *Regist) BuildUI(win fyne.Window) {
	reg.win = win

	UseSQLiteBind := binding.NewBool()
	sqlTableName := binding.NewString()
	sqliteOutFileURL := binding.NewString()
	////////////////////
	// shared widgets
	////////////////////

	// home tab buttons labels
	runButtonLBL := "Make csv table"
	outputButtonLBL := "Save CSV as"
	sqliteUpdateButtonLBL := "Update SQLite database"

	closeButton := widget.NewButtonWithIcon("Close", theme.LogoutIcon(), func() { reg.win.Close() })
	// progress bar
	//reg.progBar.Show() // to avoid resizing of the window by the progress bar
	reg.progBar.Hide()

	inputDirButton, inputDirLabel, inputDirURL := inputButton(reg)
	outFileButton, outFileLabel, outFileURL := outPutButton(reg, outputButtonLBL)
	//outFileButton, outFileLabel, outFileURL := outPutButton(reg, outputButtonLBL)

	// image logo
	pict := widget.NewCard(reg.cardTitle, reg.cardSubTitle, reg.img)

	listfiles := widget.NewCheck("List Files", func(v bool) {})
	listfiles.Checked = reg.config.ListFiles

	guessType := widget.NewCheck("Guess Dir Type", func(v bool) {})
	guessType.Checked = reg.config.GuessDirType

	dirSize := widget.NewCheck("Compute dir Size", func(v bool) {})
	// dirSize.Checked = false // set the default value to false
	dirSize.Checked = reg.config.CalcSize

	levelLab := widget.NewLabel("Level")
	level := widget.NewEntry()
	level.SetText(IntToString(reg.config.Level))
	levelEntry := container.New(layout.NewHBoxLayout(), levelLab, level)

	// Home tab status bar
	// create a label to show some information
	infoLabel := widget.NewLabel("Ready")

	////////////////
	// filters tab
	///////////////

	var selection types.RadioGroupFilters

	// create a radio button group for Names filtering
	radioGroup := widget.NewRadioGroup([]string{"Filter Names", "Path", "Path and Names"}, func(s string) {
		// do something when the option changes
		//fmt.Println("You selected", s)
		selection = updateRadioGroup(s)

	})
	selection = setRadioGroupFilters(reg, radioGroup)

	includeRegex := widget.NewCheck("Include : check to use Regex instead of string", func(v bool) {})
	//includeRegex.Checked = false // set the default value to true
	includeRegex.Checked = reg.config.IncludeRegex

	include, includeFormated := includeArea(reg)
	include.OnChanged = func(s string) {
		includeFormated = textValidation(s)
	}

	excludeRegex := widget.NewCheck("Exclude : check to use Regex instead of string", func(v bool) {})
	//excludeRegex.Checked = false // set the default value to true
	excludeRegex.Checked = reg.config.ExcludeRegex

	exclude, excludeFormated := excludeArea(reg)
	exclude.OnChanged = func(s string) {
		excludeFormated = textValidation(s)
		//log.Println(excludeFormated)
	}

	IncludeAndExclude := widget.NewCheck("Include AND Exclude (default: OR)", func(v bool) {})
	IncludeAndExclude.Checked = reg.config.IncludeAndExclude

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
		radioGroup,
		IncludeAndExclude,
		includeRegex,
		include,
		excludeRegex,
		exclude,
		dateFilter,
		container.NewGridWithColumns(2, olderthanLab, olderthan),
		container.NewGridWithColumns(2, newerthanLab, newerthan),
	)

	////////////
	// 	merge
	////////////

	// merge tab status bar
	mergeStatus := widget.NewLabel("Ready")

	// Create a string binding
	oldFileURL := binding.NewString()
	oldFileButton := getfilePath(reg.win, "Choose file to update", oldFileURL)

	// Create a string binding
	newFileURL := binding.NewString()
	newFileButton := getfilePath(reg.win, "Choose new file", newFileURL)
	mergeButton := widget.NewButton("Merge Files", func() {
		reg.progBar.Show()
		oldURL, _ := oldFileURL.Get()
		newURL, _ := newFileURL.Get()
		//merge.Merge(cleanFileURL(oldURL), cleanFileURL(newURL))
		success := merge.Merge(oldURL, newURL)
		if success {
			merge.MergeStatusDone(mergeStatus)
		} else {
			merge.MergeStatusFail(mergeStatus)
		}

		reg.progBar.Hide()

	})
	mergeContent := container.NewVBox(oldFileButton, newFileButton, mergeButton, closeButton, reg.progBar, mergeStatus)
	// Create a widget label with some help text
	helpContent := container.NewVScroll(widget.NewLabel(helpText()))

	////////////
	// run button
	////////////
	runButton := widget.NewButtonWithIcon(runButtonLBL, theme.ComputerIcon(), func() {
		go startDirAnalysis(reg, inputDirURL, outFileURL, sqliteOutFileURL, sqlTableName,
			UseSQLiteBind,
			listfiles, guessType, dirSize, includeRegex, excludeRegex, dateFilter, IncludeAndExclude,
			level, olderthan, newerthan, infoLabel,
			includeFormated, excludeFormated, selection)
	})

	////////////
	// 	SQLite
	////////////

	// sqlite tab status bar
	sqliteStatus := widget.NewLabel("Ready")
	sqliteTabLab := widget.NewLabel("SQLite Table name")
	sqliteTable := widget.NewEntry()
	sqliteTable.SetText(reg.config.SQLiteTable)
	sqlTableName.Set(reg.config.SQLiteTable)
	sqliteEntry := container.New(layout.NewVBoxLayout(), sqliteTabLab, sqliteTable)
	sqliteTable.OnChanged = func(v string) {
		sqlTableName.Set(v)
	}

	initSQLButton := widget.NewButtonWithIcon("Create SQLite table", theme.ComputerIcon(), func() {
		go makeSQLiteDB(sqliteOutFileURL, sqliteStatus)
	})

	sqliteSaveAsButton, sqliteSaveAsFileLabel := saveSQLButton(reg, "Save New SQLite database as", sqliteOutFileURL)

	sqliteOutButton, sqliteOutFileLabel := sqliteOutButton(reg, sqliteOutFileURL)

	//////////////////////////////
	//   update SQLlite button
	//////////////////////////////
	updateSQLliteButton := widget.NewButtonWithIcon(sqliteUpdateButtonLBL, theme.ComputerIcon(), func() {
		go startDirAnalysis(reg, inputDirURL, outFileURL, sqliteOutFileURL, sqlTableName,
			UseSQLiteBind,
			listfiles, guessType, dirSize, includeRegex, excludeRegex, dateFilter, IncludeAndExclude,
			level, olderthan, newerthan, infoLabel,
			includeFormated, excludeFormated, selection)
	})

	UseSQLite := widget.NewCheck("Enable SQLite database", func(v bool) {
		if v {
			updateSQLliteButton.Enable()
			updateSQLliteButton.Refresh()
		} else {
			updateSQLliteButton.Disable()
			updateSQLliteButton.Refresh()
		}
	})

	UseSQLite.Checked = reg.config.UseSQLite
	UseSQLiteBind.Set(reg.config.UseSQLite)
	// the label of the run button is changed depending if SQLite is used or not
	UseSQLite.OnChanged = func(v bool) {
		if v {
			UseSQLiteBind.Set(true)
			updateSQLliteButton.Enable()
			updateSQLliteButton.Refresh()
		} else {
			UseSQLiteBind.Set(false)
			updateSQLliteButton.Disable()
			updateSQLliteButton.Refresh()
		}
	}

	if UseSQLite.Checked {
		updateSQLliteButton.Enable()
		updateSQLliteButton.Refresh()
	} else {
		updateSQLliteButton.Disable()
		updateSQLliteButton.Refresh()
	}

	sqliteContent := container.NewVBox(UseSQLite, sqliteEntry, sqliteSaveAsButton, sqliteSaveAsFileLabel, initSQLButton, sqliteOutButton, sqliteOutFileLabel, sqliteStatus)

	//////////////////////
	// build windows tabs
	/////////////////////

	//homeContent := container.NewVBox(listfiles, guessType, dirSize, levelEntry, closeButton, pict, progBar)fyne.Window
	homeContent := container.New(layout.NewGridLayoutWithColumns(2),
		container.NewVBox(inputDirButton, inputDirLabel, outFileButton, outFileLabel, listfiles, guessType,
			dirSize, levelEntry, runButton, updateSQLliteButton, closeButton, infoLabel),
		container.NewVBox(
			pict, reg.progBar))

	homeTab := container.NewTabItem("Home", homeContent)
	filtersTab := container.NewTabItem("Filters", filtersContent)
	mergeTab := container.NewTabItem("Merge CSV", mergeContent)
	sqliteTab := container.NewTabItem("Advanced", sqliteContent)
	helpTab := container.NewTabItem("Help", helpContent)

	// build tab container
	tabs := container.NewAppTabs(homeTab, filtersTab, mergeTab, sqliteTab, helpTab)

	content := tabs

	win.SetContent(content)

	win.Content().Refresh()
}

// startDirAnalysis start a goroutine that register user settings, save them in json file
// and then start computation with cmd line engine
func startDirAnalysis(reg *Regist, inputDirURL, outFileURL, sqliteOutFileURL, sqlTableName binding.String,
	UseSQLiteBind binding.Bool,
	listfiles, guessType, dirSize, includeRegex, excludeRegex, dateFilter, IncludeAndExclude *widget.Check,
	level, olderthan, newerthan *widget.Entry,
	infoLabel *widget.Label,
	includeFormated, excludeFormated []string,
	selection types.RadioGroupFilters) {
	reg.progBar.Show()

	log.Println("Saving user settings...")
	infoLabel.Text = "Saving user settings..."
	infoLabel.Refresh()

	userSetting := reg.GetUserSettings(inputDirURL, outFileURL, sqliteOutFileURL, sqlTableName,
		UseSQLiteBind,
		listfiles, guessType, dirSize, includeRegex, excludeRegex, dateFilter, IncludeAndExclude,
		level, olderthan, newerthan,
		includeFormated, excludeFormated,
		selection)
	//log.Println(userSetting)
	reg.saveConfig(userSetting)

	log.Println("Starting directory listing...")
	infoLabel.Text = "Starting directory listing..."
	infoLabel.Refresh()

	process.Parse()

	reg.progBar.Hide()

	log.Println("Listing done !")
	infoLabel.Text = "Listing done !"
	infoLabel.Refresh()
	time.Sleep(time.Second)
	//dialog.ShowInformation("Info", "Analysis done !", reg.win) // show the info dialog

	infoLabel.Text = "Ready"
	infoLabel.Refresh()
}

// setRadioGroupFilters read and set the radiogroup configuration for filters
func setRadioGroupFilters(reg *Regist, radioGroup *widget.RadioGroup) types.RadioGroupFilters {
	var selection types.RadioGroupFilters
	savedOption := "Filter Names"

	if reg.config.FilterPath {
		savedOption = "Path"
	} else if reg.config.FilterPathName {
		savedOption = "Path and Names"
	}

	// default value from the json file
	radioGroup.SetSelected(savedOption)

	switch savedOption {
	case "Filter Names":
		selection.FilterName = true
		selection.FilterPath = false
		selection.FilterPathName = false
	case "Path":
		selection.FilterName = false
		selection.FilterPath = true
		selection.FilterPathName = false
	case "Path and Names":
		selection.FilterName = false
		selection.FilterPath = false
		selection.FilterPathName = true
	}
	//log.Println(selection)
	return selection
}

// updateRadioGroup set the user change in radiogroup selection
func updateRadioGroup(userOption string) types.RadioGroupFilters {
	var selection types.RadioGroupFilters
	switch userOption {
	case "Filter Names":
		selection.FilterName = true
		selection.FilterPath = false
		selection.FilterPathName = false
	case "Path":
		selection.FilterName = false
		selection.FilterPath = true
		selection.FilterPathName = false
	case "Path and Names":
		selection.FilterName = false
		selection.FilterPath = false
		selection.FilterPathName = true
	}
	//log.Println(selection)
	return selection
}

// inputButon return the "Choose the directory to scan"
func inputButton(reg *Regist) (*widget.Button, *widget.Label, binding.String) {
	inputDirURL := binding.NewString()
	inputDirURL.Set(reg.config.InputDir)
	inputDirStr, _ := inputDirURL.Get()
	//inputDirLabel := widget.NewLabel(inputDirStr)
	inputDirLabel := widget.NewLabelWithStyle(insertNewlines(inputDirStr, 45), fyne.TextAlignLeading, fyne.TextStyle{})
	inputDirButton := getdirPath(reg.win, "Choose the directory to scan", inputDirURL, inputDirLabel)

	return inputDirButton, inputDirLabel, inputDirURL
}

// outPutButons return the "Output file" buttons
func outPutButton(reg *Regist, outputButtonLBL string) (*widget.Button, *widget.Label, binding.String) {
	// Create a string binding
	outFileURL := binding.NewString()
	outFileURL.Set(reg.config.OutputFile)
	outFileStr, _ := outFileURL.Get()
	outFileLabel := widget.NewLabelWithStyle(insertNewlines(outFileStr, 45), fyne.TextAlignLeading, fyne.TextStyle{})
	outFileButton := getfileSave(reg.win, outputButtonLBL, outFileURL, outFileLabel) //

	return outFileButton, outFileLabel, outFileURL
}

// saveSQLButton return the "Save New SQLite database as" button
func saveSQLButton(reg *Regist, outputButtonLBL string, sqliteOutFileURL binding.String) (*widget.Button, *widget.Label) {
	// Create a string binding
	//outFileURL := binding.NewString()
	sqliteOutFileURL.Set(reg.config.OutputDB)
	sqliteOutStr, _ := sqliteOutFileURL.Get()
	sqloutFileLabel := widget.NewLabelWithStyle(insertNewlines(sqliteOutStr, 110), fyne.TextAlignLeading, fyne.TextStyle{})
	sqloutFileButton := getfileSave(reg.win, outputButtonLBL, sqliteOutFileURL, sqloutFileLabel) //

	return sqloutFileButton, sqloutFileLabel
}

// outPutButons return the "Output file" buttons
func sqliteOutButton(reg *Regist, sqliteOutFileURL binding.String) (*widget.Button, *widget.Label) {
	// Create a string binding
	sqliteOutFileURL.Set(reg.config.OutputDB)
	sqliteOutStr, _ := sqliteOutFileURL.Get()
	outFileLabel := widget.NewLabelWithStyle(insertNewlines(sqliteOutStr, 110), fyne.TextAlignLeading, fyne.TextStyle{})
	outFileButton := getDatabasePath(reg.win, "Open existing database path", sqliteOutFileURL, outFileLabel) //

	return outFileButton, outFileLabel
}

func includeArea(reg *Regist) (*widget.Entry, []string) {
	var includeFormated []string
	include := widget.NewMultiLineEntry()
	include.SetText(strSliceToString(reg.config.Include))
	includeFormated = reg.config.Include
	return include, includeFormated
}

func excludeArea(reg *Regist) (*widget.Entry, []string) {
	var excludeFormated []string
	exclude := widget.NewMultiLineEntry()
	exclude.SetText(strSliceToString(reg.config.Exclude))
	excludeFormated = reg.config.Exclude
	return exclude, excludeFormated
}

// makeSQLiteDB create a new SQLite database and display info on the status bar
func makeSQLiteDB(sqliteOutFileURL binding.String, sqliteStatus *widget.Label) {
	DBpath, _ := sqliteOutFileURL.Get()
	fmt.Println("SQLite DB path : ", DBpath)

	process.InitSQL()

	fmt.Println("Database created : " + DBpath)
	sqliteStatus.Text = "Database created : " + DBpath
	sqliteStatus.Refresh()
	time.Sleep(time.Second)
	sqliteStatus.Text = "Ready"
	sqliteStatus.Refresh()
}
