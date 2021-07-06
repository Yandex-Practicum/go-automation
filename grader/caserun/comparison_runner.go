package caserun

import "context"

type moduleComparator struct {
	runner Runner
}

var _ ModuleComparator = (*moduleComparator)(nil)

func NewComparisonRunner(runner Runner) *moduleComparator {
	return &moduleComparator{runner: runner}
}

func (r *moduleComparator) CompareModules(ctx context.Context, query ComparisonQuery) (*ComparisonReport, error) {
	originalReport, err := r.runner.Run(ctx, Query{
		ModulePath: query.OriginalModulePath,
		Suite:      query.Suite,
	})
	if err != nil {
		return nil, err
	}

	report, err := r.runner.Run(ctx, Query{
		ModulePath: query.ModulePath,
		Suite:      query.Suite,
	})
	if err != nil {
		return nil, err
	}

	return &ComparisonReport{
		OriginalSolutionReport: originalReport,
		SolutionReport:         report,
	}, nil
}
