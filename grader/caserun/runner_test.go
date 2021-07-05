package caserun_test

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/grader/caserun"
	"github.com/stretchr/testify/suite"
)

type RunnerTestSuite struct {
	suite.Suite

	tempDir string
}

func TestRunnerTestSuite(t *testing.T) {
	suite.Run(t, &RunnerTestSuite{})
}

func (s *RunnerTestSuite) SetupTest() {
	var err error
	s.tempDir, err = ioutil.TempDir("", "caserun_*")
	s.Require().NoError(err)
}

func (s *RunnerTestSuite) TearRun() {
	s.Require().NoError(os.RemoveAll(s.tempDir))
}

func (s *RunnerTestSuite) TestNoInput() {
	s.Run("PrepareCode", func() {
		s.createFile("go.mod", `
module exercise

go 1.15
`)

		s.createFile("main.go", `
package main

import "fmt"

func main() {
	fmt.Println("hello")
}
`)
	})

	s.Run("Run", func() {
		runner := caserun.NewRunner()
		report, err := runner.Run(context.Background(), caserun.Query{
			ModulePath: s.tempDir,
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
		s.EqualValues("hello\n", caseReport.UserOutput)
		s.Less(caseReport.TimeUsed, caseReport.TimeLimitMilli)
	})
}

func (s *RunnerTestSuite) createFile(name, content string) {
	s.Require().NoError(ioutil.WriteFile(
		path.Join(s.tempDir, name),
		[]byte(content),
		os.ModePerm,
	))
}
