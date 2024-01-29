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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

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
			label.SetText(cleanFileURL(file.URI().String()))
			//path = file.URI().String()
			url.Set(cleanFileURL(file.URI().String()))
		}, window)
	}
	// Créer un bouton qui déclenche la fonction de choix de fichier
	button := widget.NewButton(buttonlabel, chooseFile)

	return container.NewVBox(button, label)
}

// getfileSave create a file button and stores the file path entered by the user and refresh the path label
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
			url.Set(cleanFileURL(file.URI().String()))
			outFileLabel.Text = insertNewlines(cleanFileURL(file.URI().String()), 45)
			outFileLabel.Refresh()
		}, window)
	}
	// Créer un bouton qui déclenche la fonction de choix de fichier
	button := widget.NewButton(buttonlabel, chooseFile)

	return button
}

// getdirPath create a file button and stores the dir path using databinding and refresh the path label
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
			outDirLabel.SetText(insertNewlines(cleanFileURL(dir.Path()), 45))
			url.Set(cleanFileURL(dir.String()))
			outDirLabel.Refresh()
		}, window)
	}
	// Créer un bouton qui déclenche la fonction de choix de fichier
	button := widget.NewButton(buttonlabel, chooseDir)

	return button
}

// getDatabasePath stores the SQLIte database path
func getDatabasePath(window fyne.Window, buttonlabel string, url binding.String, outFileLabel *widget.Label) *widget.Button {
	//var path string // file path
	// Créer un label pour afficher le chemin du fichier
	//label := widget.NewLabel("")

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
			// exists, err := storage.Exists(file)
			// if err != nil {
			// 	// handle error
			// }
			// if !exists {
			// 	// Fermer le fichier
			// 	file.Close()
			// }
			// Afficher le chemin du fichier dans le label
			//label.SetText(file.URI().String())
			//path = file.URI().String()
			url.Set(cleanFileURL(file.URI().String()))
			outFileLabel.Text = insertNewlines(cleanFileURL(file.URI().String()), 100)
			outFileLabel.Refresh()
		}, window)
	}
	// Créer un bouton qui déclenche la fonction de choix de fichier
	button := widget.NewButton(buttonlabel, chooseFile)

	return button
}
