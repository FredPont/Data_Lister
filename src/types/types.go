package types

type Conf struct {
	InputDir     string   `json:"InputDir"`
	OutputFile   string   `json:"OutputFile"`
	ListFiles    bool     `json:"ListFiles"`
	ListDir      bool     `json:"ListDir"`
	GuessDirType bool     `json:"GuessDirType"`
	CalcSize     bool     `json:"CalcSize"`
	Level        int      `json:"Level"`
	IncludeRegex bool     `json:"IncludeRegex"`
	Include      []string `json:"Include"`
	ExcludeRegex bool     `json:"ExcludeRegex"`
	Exclude      []string `json:"Exclude"`
	OlderThan    string   `json:"OlderThan"`
	NewerThan    string   `json:"NewerThan"`
}

type DirSignature struct {
	Content        []string `json:"content"`
	ScoreThreshold float64  `json:"scoreThreshold"`
}

type DirMatch struct {
	IsMatch bool
	Label   string
	Score   float64
}
