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
