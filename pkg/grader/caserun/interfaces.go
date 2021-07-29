package caserun

import "context"

type ModuleComparator interface {
	CompareModules(ctx context.Context, query ComparisonQuery) (*ComparisonReport, error)
}

type Runner interface {
	Run(ctx context.Context, query Query) (*SuiteReport, error)
}
