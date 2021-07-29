package caserun

type ComparisonQuery struct {
	OriginalModulePath string `json:"original_module_path"`
	ModulePath         string `json:"module_path"`
	Suite              Suite  `json:"suite"`
}

type Query struct {
	ModulePath string `json:"module_path"`
	Suite      Suite  `json:"suite"`
}
