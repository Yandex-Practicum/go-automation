package caserun_test

import (
	"context"
	"testing"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/grader/caserun"
	"github.com/stretchr/testify/suite"
)

type RunnerTestSuite struct {
	Suite
}

func TestRunnerTestSuite(t *testing.T) {
	suite.Run(t, &RunnerTestSuite{})
}

func (s *RunnerTestSuite) TestNoInput() {
	s.Run("PrepareCode", func() {
		s.CreateMod(s.moduleDir)
		s.CreateMain(`
package main

import "fmt"

func main() {
	fmt.Println("hello")
}
`, s.moduleDir)
	})

	s.Run("Run", func() {
		runner := caserun.NewRunner()
		report, err := runner.Run(context.Background(), caserun.Query{
			ModulePath: s.moduleDir,
			Suite: caserun.Suite{
				ID: "Suite",
				Cases: []caserun.Case{
					{
						ID:             "id",
						Tag:            "tag",
						TimeLimitMilli: 1000,
					},
				},
			},
		})
		s.Require().NoError(err)

		s.Require().NotNil(report)
		s.Require().Len(report.Cases, 1)

		caseReport := report.Cases[0]
		s.EqualValues("id", caseReport.ID)
		s.EqualValues("tag", caseReport.Tag)
		s.EqualValues(1000, caseReport.TimeLimitMilli)
		s.EqualValues("hello\n", caseReport.Stdout)
		s.Less(caseReport.TimeUsedMilli, caseReport.TimeLimitMilli)
	})
}

func (s *RunnerTestSuite) TestWithUserInput() {
	s.Run("PrepareCode", func() {
		s.CreateMod(s.moduleDir)
		s.CreateMain(`
package main

import (
	"fmt"
)

func main() {
	var input string
    fmt.Scanln(&input)
	fmt.Println(input + "_suffix")
}
`, s.moduleDir)
	})

	s.Run("Run", func() {
		runner := caserun.NewRunner()
		report, err := runner.Run(context.Background(), caserun.Query{
			ModulePath: s.moduleDir,
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

		s.Require().NotNil(report)
		s.Require().Len(report.Cases, 1)

		caseReport := report.Cases[0]
		s.EqualValues("id", caseReport.ID)
		s.EqualValues("tag", caseReport.Tag)
		s.EqualValues(1000, caseReport.TimeLimitMilli)
		s.EqualValues("input_suffix\n", caseReport.Stdout)
		s.Less(caseReport.TimeUsedMilli, caseReport.TimeLimitMilli)
	})
}
