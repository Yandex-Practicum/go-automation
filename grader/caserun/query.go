package caserun

type ComparisonQuery struct {
	OriginalModulePath string `json:"originalModulePath"`
	ModulePath         string `json:"modulePath"`
	Suite              Suite  `json:"suite"`
}

type Query struct {
	ModulePath string `json:"modulePath"`
	Suite      Suite  `json:"suite"`
}
