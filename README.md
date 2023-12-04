# DATA Lister
 DATA Lister is a software to list directories/files and produce a TSV table for data management.

  ![Output Table](src/images/table.png)
work in progress...

## Key characteristics
- List directories and optionally files.
- TSV output
- Tunable dir level
- Compute dir size
- Try to guess dir content (level must be set to dir +1 to allow the dir content analysis)
- Include/Exclude filter list by string or regex
- Filter by date
- Custom (pre-filed) columns

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
    "CalcSize": true,
    "Level": 4,
    "IncludeRegex": true,
    "Include": [
      ".*"
    ],
    "ExcludeRegex": true,
    "Exclude": [
        "image"
    ],
    "DateFilter": false,
    "OlderThan": "2023-12-12",
    "NewerThan": "2022-12-12"
}
```
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
ColumnName	DefaultValues
SampleType	Cells
Project_ID	Project_1
RelatedProject	Project_2
content	MyExperiments
Delete_Date	2028-01-01
```
- Start the software using the precompiled binaries for Linux, Mac or Windows
