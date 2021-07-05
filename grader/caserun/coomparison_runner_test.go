package caserun_test

import (
	"context"
	"testing"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/grader/caserun"
	"github.com/stretchr/testify/suite"
)

type ComparisonRunnerTestSuite struct {
	Suite
}

func TestComparisonRunnerTestSuite(t *testing.T) {
	suite.Run(t, &ComparisonRunnerTestSuite{})
}

func (s *ComparisonRunnerTestSuite) TestNotMatchingSolution() {
	s.Run("SetUpOriginalSolution", func() {
		s.CreateMod(s.originalModuleDir)
		s.CreateMain(`
package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello")
}
`, s.originalModuleDir)
	})

	s.Run("SetUpUserSolution", func() {
		s.CreateMod(s.moduleDir)
		s.CreateMain(`
package main

import (
	"fmt"
)

func main() {
	fmt.Println("buy")
}
`, s.moduleDir)
	})

	s.Run("MakeComparisonReport", func() {
		runner := caserun.NewComparisonRunner(caserun.NewRunner())
		report, err := runner.Run(context.Background(), caserun.ComparisonQuery{
			OriginalModulePath: s.originalModuleDir,
			ModulePath:         s.moduleDir,
			Suite: caserun.Suite{
				ID: "Suite",
				Cases: []caserun.Case{
					{
						ID:             "id",
						Tag:            "tag",
						Input:          "input",
						TimeLimitMilli: 1000,
					},
				},
			},
		})
		s.Require().NoError(err)

		s.EqualValues("hello\n", report.OriginalSolutionReport.Cases[0].UserOutput)
		s.EqualValues("buy\n", report.SolutionReport.Cases[0].UserOutput)
	})
}
