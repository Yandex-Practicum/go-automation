package caserun

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/grader"
	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/grader/commandrun"
	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/grader/compilation"
)

type runner struct{}

var _ Runner = (*runner)(nil)

func NewRunner() Runner {
	return &runner{}
}

func (r *runner) Run(ctx context.Context, query Query) (*SuiteReport, error) {
	solutionDir, err := r.makeSolutionDir(ctx, query.Suite.ID)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(solutionDir)

	binaryPath, err := r.buildBinary(ctx, solutionDir, query)
	if err != nil {
		return nil, err
	}

	caseReports := make([]CaseReport, 0, len(query.Suite.Cases))
	for _, c := range query.Suite.Cases {
		caseReport, err := r.runCase(ctx, binaryPath, c)
		if err != nil {
			return nil, err
		}

		caseReports = append(caseReports, *caseReport)
	}

	return &SuiteReport{
		Cases: caseReports,
	}, nil
}

func (r *runner) runCase(ctx context.Context, binaryPath string, c Case) (*CaseReport, error) {
	testRunner := commandrun.NewRunner(c.TimeLimitMilli.Duration(), binaryPath)

	runResult, err := testRunner.Run(ctx, commandrun.RunOptions{
		Stdin: c.Input,
	})
	if err != nil {
		return nil, err
	}

	return &CaseReport{
		ID:             c.ID,
		Tag:            c.Tag,
		TimeLimitMilli: c.TimeLimitMilli,
		TimeUsedMilli:  grader.TimeMilli(runResult.Duration / time.Millisecond),
		Stdin:          c.Input,
		Stdout:         runResult.Stdout,
		Stderr:         runResult.Stderr,
		ExitCode:       runResult.ExitCode,
	}, nil
}

func (r *runner) buildBinary(ctx context.Context, solutionDir string, query Query) (string, error) {
	binaryPath := path.Join(solutionDir, "bin")
	compiler := compilation.NewCompiler()

	if err := compiler.CompilePackage(ctx, compilation.Query{
		ModulePath: query.ModulePath,
		BinaryPath: binaryPath,
		Timeout:    time.Second * 10,
	}); err != nil {
		return "", err
	}

	return binaryPath, nil
}

func (r *runner) makeSolutionDir(ctx context.Context, suiteID SuiteID) (string, error) {
	return ioutil.TempDir("", fmt.Sprintf("suite_%s_*", suiteID))
}
