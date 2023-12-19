 <!--![Output Table](src/images/logo.png) -->


<p align="center">
<img src="src/images/logo.png" alt="drawing" width="250" height="250" />
</p>

# DATA Lister
 DATA Lister is a software to list directories/files and produce a TSV table for data management.

![GUI](src/images/gui1.png)

![GUI](src/images/gui2.png)

![GUI](src/images/gui3.png)
work in progress...

  ![Output Table](src/images/table.png)


Caution : set the directory level to a reasonable value before starting DATA Lister on a large file system to avoid producing a huge result table.

## Key characteristics
- List directories and optionally files.
- TSV table output
- Tunable dir level
- Try to guess dir content based on footprint (level must be set to dir +1 to allow the dir content analysis)
- Unlimited number of customizable dir footprint
- Include/Exclude filter list by unlimited number of string or regex
- Filter by date
- Unlimited number of custom (pre-filed) columns
- Compute directory size (slow on terabytes of data because 100% of the files are parsed)
- Merge tool to update old table with the new rows from a new analysis

## Installation

No installation required the code is statically compiled.

- Download the zip file from the "<>Code" green button and unzip it 
- or `git clone https://github.com/FredPont/Data_Lister.git`

## Quick start

- Edit config/settings.json to set root directory and options
```json
{
    "InputDir": "test",
    "OutputFile": "results/table.csv",
    "ListFiles": false,
    "GuessDirType": true,
    "CalcSize": false,
    "Level": 4,
    "IncludeRegex": true,
    "Include": [
	      ".*", "level", "seq"	
    ]
    "ExcludeRegex": true,
    "Exclude": [
        "image"
    ],
    "DateFilter": false,
    "OlderThan": "2023-12-12",
    "NewerThan": "2022-12-12"
}
```
Note : backslashes must be escaped in regex. 

Example : to exclude names starting with a dot use "^\\\\..+"

- Edit config/DirSignatures.json to set the directory patterns (strings, no regex)
```json
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
```
- Edit config/columns.tsv to add custom columns and their optional default values
```tsv
ColumnName	DefaultValueswork in progress...
SampleType	Cells
Project_ID	Project_1
RelatedProject	Project_2
content	MyExperiments
Delete_Date	2028-01-01
```
- Start the software using the precompiled binaries for Linux, Mac or Windows
```bash
Usage :
  -c	Start DataLister directories analysis in command line.
  -g	Start DataLister directories analysis in graphic mode.
  -m	Start DataLister merging tool.
  -i string
    	New result file path. Only new files/dir are added to the old file
  -o string
    	Old result file path. 

      Examples :

      Start the analysis of the directories in command line (-c):
      ./Linux_DataLister.bin -c

      Start the analysis of the directories in graphic mode (-g):
      ./Linux_DataLister.bin -g

      To add new data from newfile to oldfile :
      ./Linux_DataLister.bin -m -o oldfile.csv -i newfile.csv
```