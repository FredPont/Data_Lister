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
	"fyne.io/fyne/v2/canvas"
)

// picture display an image in the main window
func picture() fyne.CanvasObject {
	var img = canvas.NewImageFromFile("src/ui/logo.png")

	img.SetMinSize(fyne.Size{Width: 200, Height: 150})
	img.FillMode = canvas.ImageFillContain
	//img.FillMode = canvas.ImageFillOriginal
	return img
}
