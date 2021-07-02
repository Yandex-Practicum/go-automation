package commandrun_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/grader/commandrun"
	"github.com/stretchr/testify/suite"
)

type RunnerTestSuite struct {
	suite.Suite

	tempDirPath string
}

func TestRunnerTestSuite(t *testing.T) {
	suite.Run(t, &RunnerTestSuite{})
}

func (s *RunnerTestSuite) TestStdout() {
	runner := commandrun.NewRunner(time.Second, "echo", "hello")

	result, err := runner.Run(context.Background())
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.EqualValues("hello\n", result.Stdout)
}

func (s *RunnerTestSuite) TestStderr() {
	runner := commandrun.NewRunner(time.Second, "logger", "-s", "hello")

	result, err := runner.Run(context.Background())
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.True(strings.HasSuffix(result.Stderr, "hello\n"), result.Stderr)
}

func (s *RunnerTestSuite) TestTimeout() {
	runner := commandrun.NewRunner(time.Millisecond, "sleep", "1")

	result, err := runner.Run(context.Background())
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Less(result.Duration, time.Millisecond*4)
}
