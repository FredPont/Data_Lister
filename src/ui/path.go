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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// getfilePath create a file button and stores the file path using databinding
func getfilePath(window fyne.Window, buttonlabel string, url binding.String) *fyne.Container {

	// label to display file path
	label := widget.NewLabel("")

	// file choosing function
	chooseFile := func() {
		// dialog to open a file
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil {
				// show error on pop up
				dialog.ShowError(err, window)
				return
			}
			if file == nil {
				// do nothing when no file is selected
				return
			}
			// close the file
			file.Close()
			// display file path in label
			label.SetText(cleanFileURI(file.URI()))
			//path = file.URI().String()
			url.Set(cleanFileURI(file.URI()))
		}, window)
	}
	// button to trigger file pickup
	button := widget.NewButtonWithIcon(buttonlabel, theme.FileIcon(), chooseFile)

	return container.NewVBox(button, label)
}

// getfileSave create a file button and stores the file path entered by the user and refresh the path label
func getfileSave(window fyne.Window, buttonlabel string, url binding.String, outFileLabel *widget.Label) *widget.Button {

	// file choosing function
	chooseFile := func() {
		// dialog to open a file
		dialog.ShowFileSave(func(file fyne.URIWriteCloser, err error) {
			if err != nil {
				// show error on pop up
				dialog.ShowError(err, window)
				return
			}
			if file == nil {
				// do nothing when no file is selected
				return
			}

			// close the file
			file.Close()

			url.Set(cleanFileURI(file.URI()))
			outFileLabel.Text = insertNewlines(cleanFileURI(file.URI()), 45)
			outFileLabel.Refresh()
		}, window)
	}
	// button to trigger file pickup
	button := widget.NewButtonWithIcon(buttonlabel, theme.FileIcon(), chooseFile)

	return button
}

// getdirPath create a file button and stores the dir path using databinding and refresh the path label
func getdirPath(window fyne.Window, buttonlabel string, url binding.String, outDirLabel *widget.Label) *widget.Button {

	// dir choosing function
	chooseDir := func() {
		// dialog to open a dir
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil {
				// show error on pop up
				dialog.ShowError(err, window)
				return
			}
			if dir == nil {
				// do nothing when no file is selected
				return
			}

			outDirLabel.SetText(insertNewlines(cleanDirURI(dir), 45))
			url.Set(cleanDirURI(dir))
			outDirLabel.Refresh()
		}, window)
	}
	// button to trigger dir pickup
	button := widget.NewButtonWithIcon(buttonlabel, theme.FolderOpenIcon(), chooseDir)

	return button
}

// getDatabasePath stores the SQLIte database path
func getDatabasePath(window fyne.Window, buttonlabel string, url binding.String, outFileLabel *widget.Label) *widget.Button {

	// file choosing function
	chooseFile := func() {
		// dialog to open a file
		dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
			if err != nil {
				// show error on pop up
				dialog.ShowError(err, window)
				return
			}
			if file == nil {
				// do nothing when no file is selected
				return
			}

			url.Set(cleanFileURI(file.URI()))
			outFileLabel.Text = insertNewlines(cleanFileURI(file.URI()), 100)
			outFileLabel.Refresh()
		}, window)
	}
	// button to trigger file pickup
	button := widget.NewButtonWithIcon(buttonlabel, theme.FileIcon(), chooseFile)

	return button
}
