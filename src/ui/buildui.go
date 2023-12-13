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

	regist := &Regist{img: im}
	return regist
}

// BuildUI creates the main window of our application
func (reg *Regist) BuildUI(w fyne.Window) {
	reg.win = w

	//---------
	// home tab

	inputDirURL := binding.NewString()
	inputDirButton := getdirPath(reg.win, "Choose the directory to scan", inputDirURL)

	//oldFileURL.Set(oldfile)
	// Create a string binding
	outFileURL := binding.NewString()
	outFileButton := getfileSave(reg.win, "Output file", outFileURL)

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

	closeButton := widget.NewButtonWithIcon("Close", theme.LogoutIcon(), func() { reg.win.Close() })

	progBar := widget.NewProgressBarInfinite()
	progBar.Hide()

	//homeContent := container.NewVBox(listfiles, guessType, dirSize, levelEntry, closeButton, pict, progBar)
	homeContent := container.New(layout.NewGridLayoutWithColumns(2),
		container.NewVBox(inputDirButton, outFileButton, listfiles, guessType, dirSize, levelEntry, closeButton, progBar),
		pict)

	//-------------
	// filters tab
	includeRegex := widget.NewCheck("Include Regex", func(v bool) {})
	includeRegex.Checked = false // set the default value to true

	includeFormated := ""
	include := widget.NewMultiLineEntry()
	include.OnChanged = func(s string) {
		includeFormated = commentValidation(s)
		log.Println(includeFormated)
	}

	excludeRegex := widget.NewCheck("Exclude Regex", func(v bool) {})
	excludeRegex.Checked = false // set the default value to true
	excludeFormated := ""
	exclude := widget.NewMultiLineEntry()
	exclude.OnChanged = func(s string) {
		excludeFormated = commentValidation(s)
		log.Println(includeFormated)
		log.Println(excludeFormated)
	}

	dateFilter := widget.NewCheck("Date Filter", func(v bool) {})
	dateFilter.Checked = false // set the default value to true

	olderthan := widget.NewEntry()
	olderthan.SetPlaceHolder("3023-12-12")

	newerthan := widget.NewEntry()
	newerthan.SetPlaceHolder("1922-12-12")

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

	//oldFileURL.Set(oldfile)
	// Create a string binding
	newFileURL := binding.NewString()
	newFileButton := getfilePath(reg.win, "Choose new file", newFileURL)
	mergeButton := widget.NewButton("Merge Files", func() {
		progBar.Show()
		oldURL, _ := oldFileURL.Get()
		newURL, _ := newFileURL.Get()
		time.Sleep(3 * time.Second) // pauses the execution for 3 seconds
		log.Println(oldURL, newURL)
		progBar.Hide()
	})
	mergeContent := container.NewVBox(oldFileButton, newFileButton, mergeButton, progBar)
	// Create a widget label with some help text
	helpContent := widget.NewLabel(helpText())

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
func getfilePath(window fyne.Window, buttonlabel string, url binding.String) *fyne.Container {
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

// getfileSave create a file button and stores the file path entered by the user
func getfileSave(window fyne.Window, buttonlabel string, url binding.String) *fyne.Container {
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
func getdirPath(window fyne.Window, buttonlabel string, url binding.String) *fyne.Container {
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

// getEntry create an entry widget and stores the entry
func getEntry() {

}

func commentValidation(cmt string) string {
	cmt = strings.ReplaceAll(cmt, "\n", " ")  // remove all newline
	return strings.ReplaceAll(cmt, "\r", " ") // remove all carriage return
}
