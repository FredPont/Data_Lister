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
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
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
	level.SetPlaceHolder("3")
	levelEntry := container.New(layout.NewHBoxLayout(), levelLab, level)

	// buttons
	closeButton := widget.NewButtonWithIcon("Close", theme.LogoutIcon(), func() { reg.win.Close() })
	runButton := widget.NewButtonWithIcon("Run", theme.ComputerIcon(), func() {
		reg.saveConfig(userSetting)
		//reg.win.Close()
	})

	progBar := widget.NewProgressBarInfinite()
	progBar.Hide()

	//homeContent := container.NewVBox(listfiles, guessType, dirSize, levelEntry, closeButton, pict, progBar)
	homeContent := container.New(layout.NewGridLayoutWithColumns(2),
		container.NewVBox(inputDirButton, inputDirLabel, outFileButton, outFileLabel, listfiles, guessType, dirSize, levelEntry, runButton, closeButton, progBar),
		pict)

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
	oldFileButton := getfilePath2(reg.win, "Choose file to update", oldFileURL)
	//oldFileButton

	//oldFileURL.Set(oldfile)
	// Create a string binding
	newFileURL := binding.NewString()
	newFileButton := getfilePath2(reg.win, "Choose new file", newFileURL)
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

// picture display an image in the main window
func picture() fyne.CanvasObject {
	var img = canvas.NewImageFromFile("src/ui/logo.png")

	img.SetMinSize(fyne.Size{Width: 330, Height: 220})
	img.FillMode = canvas.ImageFillContain
	return img
}

// helpText() return the text hold by the help tab
func helpText() string {
	text := `
	Usage :
	-c	Start DataLister directories analysis in command line.

	-m	Start DataLister merging tool.
	-i string
	New result file path. Only new files/dir are added to the old file
	-o string
    	Old result file path. 

	Examples :

	Start the analysis of the directories in command line (-c):
	./Linux_DataLister.bin -c

	To add new data from newfile to oldfile :
	./Linux_DataLister.bin -m -o oldfile.csv -i newfile.csv
	`

	return text

}

// getfilePath create a file button and stores the file path using databinding
func getfilePath2(window fyne.Window, buttonlabel string, url binding.String) *fyne.Container {
	//var path string // file path
	// Créer un label pour afficher le chemin du fichier
	label := widget.NewLabel("")

	// Créer une fonction de choix de fichier
	chooseFile := func() {
		// Ouvrir une boîte de dialogue pour sélectionner un fichier
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil {
				// Afficher une erreur si nécessaire
				dialog.ShowError(err, window)
				return
			}
			if file == nil {
				// Ne rien faire si aucun fichier n'est sélectionné
				return
			}
			// Fermer le fichier
			file.Close()
			// Afficher le chemin du fichier dans le label
			label.SetText(file.URI().String())
			//path = file.URI().String()
			url.Set(file.URI().String())
		}, window)
	}
	// Créer un bouton qui déclenche la fonction de choix de fichier
	button := widget.NewButton(buttonlabel, chooseFile)

	// Ajouter le bouton et le label à la fenêtre
	//	window.SetContent(container.NewVBox(button, label))

	return container.NewVBox(button, label)
}

func getfilePath(window fyne.Window, buttonlabel string, url binding.String, outFileLabel *widget.Label) *widget.Button {
	// Créer une fonction de choix de fichier
	chooseFile := func() {
		// Ouvrir une boîte de dialogue pour sélectionner un fichier
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil {
				// Afficher une erreur si nécessaire
				dialog.ShowError(err, window)
				return
			}
			if file == nil {
				// Ne rien faire si aucun fichier n'est sélectionné
				return
			}
			// Fermer le fichier
			file.Close()
			// Afficher le chemin du fichier dans le label
			//label.SetText(file.URI().String())
			//path = file.URI().String()
			url.Set(file.URI().String())
			outFileLabel.Text = insertNewlines(file.URI().String(), 45)
			outFileLabel.Refresh()
		}, window)
	}
	// Créer un bouton qui déclenche la fonction de choix de fichier
	button := widget.NewButton(buttonlabel, chooseFile)

	return button
}

// getfileSave create a file button and stores the file path entered by the user
func getfileSave(window fyne.Window, buttonlabel string, url binding.String, outFileLabel *widget.Label) *widget.Button {
	//var path string // file path
	// Créer un label pour afficher le chemin du fichier
	//label := widget.NewLabel("")

	// Créer une fonction de choix de fichier
	chooseFile := func() {
		// Ouvrir une boîte de dialogue pour sélectionner un fichier
		dialog.ShowFileSave(func(file fyne.URIWriteCloser, err error) {
			if err != nil {
				// Afficher une erreur si nécessaire
				dialog.ShowError(err, window)
				return
			}
			if file == nil {
				// Ne rien faire si aucun fichier n'est sélectionné
				return
			}
			// Fermer le fichier
			file.Close()
			// Afficher le chemin du fichier dans le label
			//label.SetText(file.URI().String())
			//path = file.URI().String()
			url.Set(file.URI().String())
			outFileLabel.Text = insertNewlines(file.URI().String(), 45)
			outFileLabel.Refresh()
		}, window)
	}
	// Créer un bouton qui déclenche la fonction de choix de fichier
	button := widget.NewButton(buttonlabel, chooseFile)

	// Ajouter le bouton et le label à la fenêtre
	//	window.SetContent(container.NewVBox(button, label))

	//return container.NewVBox(button, label)
	return button
}

func getfileSave2(window fyne.Window, buttonlabel string, url binding.String) *fyne.Container {
	//var path string // file path
	// Créer un label pour afficher le chemin du fichier
	label := widget.NewLabel("")

	// Créer une fonction de choix de fichier
	chooseFile := func() {
		// Ouvrir une boîte de dialogue pour sélectionner un fichier
		dialog.ShowFileSave(func(file fyne.URIWriteCloser, err error) {
			if err != nil {
				// Afficher une erreur si nécessaire
				dialog.ShowError(err, window)
				return
			}
			if file == nil {
				// Ne rien faire si aucun fichier n'est sélectionné
				return
			}
			// Fermer le fichier
			file.Close()
			// Afficher le chemin du fichier dans le label
			label.SetText(file.URI().String())
			//path = file.URI().String()
			url.Set(file.URI().String())
		}, window)
	}
	// Créer un bouton qui déclenche la fonction de choix de fichier
	button := widget.NewButton(buttonlabel, chooseFile)

	// Ajouter le bouton et le label à la fenêtre
	//	window.SetContent(container.NewVBox(button, label))

	return container.NewVBox(button, label)
}

// getdirPath create a file button and stores the dir path using databinding
func getdirPath(window fyne.Window, buttonlabel string, url binding.String, outDirLabel *widget.Label) *widget.Button {
	//var path string // file path
	// Créer un label pour afficher le chemin du fichier
	//label := widget.NewLabel("")

	// Créer une fonction de choix de fichier
	chooseDir := func() {
		// Ouvrir une boîte de dialogue pour sélectionner un fichier
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil {
				// Afficher une erreur si nécessaire
				dialog.ShowError(err, window)
				return
			}
			if dir == nil {
				// Ne rien faire si aucun fichier n'est sélectionné
				return
			}

			// Afficher le chemin du fichier dans le label
			//label.SetText(dir.Path())
			//path = file.URI().String()
			outDirLabel.SetText(insertNewlines(dir.Path(), 45))
			url.Set(dir.String())
			outDirLabel.Refresh()
		}, window)
	}
	// Créer un bouton qui déclenche la fonction de choix de fichier
	button := widget.NewButton(buttonlabel, chooseDir)

	// Ajouter le bouton et le label à la fenêtre
	//	window.SetContent(container.NewVBox(button, label))

	return button
}

// getdirPath create a file button and stores the dir path using databinding
func getdirPath2(window fyne.Window, buttonlabel string, url binding.String) *fyne.Container {
	//var path string // file path
	// Créer un label pour afficher le chemin du fichier
	label := widget.NewLabel("")

	// Créer une fonction de choix de fichier
	chooseDir := func() {
		// Ouvrir une boîte de dialogue pour sélectionner un fichier
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil {
				// Afficher une erreur si nécessaire
				dialog.ShowError(err, window)
				return
			}
			if dir == nil {
				// Ne rien faire si aucun fichier n'est sélectionné
				return
			}

			// Afficher le chemin du fichier dans le label
			label.SetText(dir.Path())
			//path = file.URI().String()
			url.Set(dir.String())
		}, window)
	}
	// Créer un bouton qui déclenche la fonction de choix de fichier
	button := widget.NewButton(buttonlabel, chooseDir)

	// Ajouter le bouton et le label à la fenêtre
	//	window.SetContent(container.NewVBox(button, label))

	return container.NewVBox(button, label)
}

func textValidation(cmt string) []string {
	//cmt = strings.ReplaceAll(cmt, "\n", " ")  // remove all newline
	cmt = strings.ReplaceAll(cmt, "\r", " ") // remove all carriage return
	return strings.Split(cmt, "\n")

}

// cleanFileURL removes "file://" at the beginning of URL
func cleanFileURL(url string) string {
	return strings.TrimPrefix(url, "file://") // t is "lang"
}

// dateValidator control the date format
func dateValidator(text string) error {
	const dateFormat = "2006-01-02"
	_, err := time.Parse(dateFormat, text)
	return err
}

// saveConfig export the user setting to the config/settingsjson file
func (reg *Regist) saveConfig(userSetting types.Conf) {

}

// formatURL insert newline every x char to split the url in more than one line
func formatURL(url string) {
	//x := 50
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
