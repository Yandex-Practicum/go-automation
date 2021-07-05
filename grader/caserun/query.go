package caserun

type ComparisonQuery struct {
	OriginalModulePath string
	ModulePath         string
	Suite              Suite
}

type Query struct {
	ModulePath string
	Suite      Suite
}
