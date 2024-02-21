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
	-g	Start DataLister directories analysis in graphic mode.
	-m	Start DataLister merging tool.
	-i string
	New TSV file path. Only new files/dir are added to the old file
	-o string
    	Old TSV file path. 
	-s	Create a new SQLite database. Example : DataLister -s
	Examples :

	Start the analysis of the directories in command line (-c):
	./Linux_DataLister.bin -c

	To add new data from newfile to oldfile :
	./Linux_DataLister.bin -m -o oldfile.csv -i newfile.csv

	if "UseSQLite": true in the config/settings.json file, then
	./Linux_DataLister.bin -c
	will update the SQLite database indicated in "OutputSQLFile"


	Quick start :
	
	- Edit config/settings.json to set root directory and options

	{
		"InputDir": "test",
		"OutputCSVFile": "results/table.csv",
		"OutputSQLFile": "test/testDB.db",
		"ListFiles": false,
		"GuessDirType": false,
		"CalcSize": true,
		"Level": 3,
		"filterName": false,
		"filterPath": false,
		"filterPathName": true,
		"IncludeRegex": false,
		"Include": [
			""
		],
		"IncludeAndExclude": false,
		"ExcludeRegex": false,
		"Exclude": [
			""
		],
		"OlderThan": "3023-12-12",
		"NewerThan": "1922-12-12",
		"DateFilter": false,
		"UseSQLite": false,
		"SQLiteTable": "table1",
		"CompiledIncludeRegex": null,
		"CompiledExcludeRegex": null
	}

	Use absolute path in "InputDir", "OutputCSVFile" or "OutputSQLFile".
	Note : for the command line version, backslashes must be escaped 
	in regex in the settings.json file (this is not necessary in the GUI). 

	Example : to exclude names starting with a dot use "^\\..+"

	- The filter priority is Date > Include > Exclude
	- If more than one string/regex is given they are cumulated (reg1 OR reg2)
	- If Include and Exclude are used simultaneously, they are cumulated (Include OR Exclude)
	  if "Include AND Exclude" is not checked

	- Edit config/DirSignatures.json to set the directory patterns (strings, no regex)

	{
		"Software": {
					"content": [".go", ".git", ".DLL", ".dll", ".r", ".jl", ".pl"],
					"scoreThreshold": 0.2
		},
		"Fasta": {
				"content": [".fasta", ".FASTA", ".fasta.gz"],
				"scoreThreshold": 0.8
		}
	}

	- Edit config/columns.tsv to add custom columns and their optional default values
	tsv
	ColumnName	DefaultValueswork in progress...
	SampleType	Cells
	Project_ID	Project_1
	RelatedProject	Project_2
	content	MyExperiments
	Delete_Date	2028-01-01
	`

	return text

}
