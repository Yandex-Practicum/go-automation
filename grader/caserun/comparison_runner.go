package caserun

import "context"

type ComparisonRunner interface {
	Run(ctx context.Context, query ComparisonQuery) (*ComparisonReport, error)
}

type comparisonRunner struct {
	runner Runner
}

var _ ComparisonRunner = (*comparisonRunner)(nil)

func NewComparisonRunner(runner Runner) *comparisonRunner {
	return &comparisonRunner{runner: runner}
}

func (r *comparisonRunner) Run(ctx context.Context, query ComparisonQuery) (*ComparisonReport, error) {
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
