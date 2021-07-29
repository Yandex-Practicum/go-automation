package compilation_test

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/grader/commandrun"
	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/grader/compilation"
	"github.com/stretchr/testify/suite"
)

type CompilerTestSuite struct {
	suite.Suite

	tempDir string
}

func TestCompilerTestSuite(t *testing.T) {
	suite.Run(t, &CompilerTestSuite{})
}

func (s *CompilerTestSuite) SetupTest() {
	var err error
	s.tempDir, err = ioutil.TempDir("", "compilation_*")
	s.Require().NoError(err)
}

func (s *CompilerTestSuite) TearDownTest() {
	s.Require().NoError(os.RemoveAll(s.tempDir))
}

func (s *CompilerTestSuite) TestCompileAndRun() {
	s.Run("Compile", func() {
		s.createFile("go.mod", `
module exercise

go 1.15
`)

		s.createFile("main.go", `
package main

func main() {
	println("hello")
}
`)

		compiler := compilation.NewCompiler()
		s.Require().NoError(compiler.CompilePackage(context.Background(), compilation.Query{
			ModulePath: s.tempDir,
			Timeout:    time.Minute,
		}))
	})

	s.Run("Run", func() {
		runner := commandrun.NewRunner(time.Second, path.Join(s.tempDir, "exercise"))
		runResult, err := runner.Run(context.Background(), commandrun.RunOptions{})
		s.Require().NoError(err)

		s.EqualValues("hello\n", runResult.Stderr)
	})
}

func (s *CompilerTestSuite) createFile(name, content string) {
	s.Require().NoError(ioutil.WriteFile(
		path.Join(s.tempDir, name),
		[]byte(content),
		os.ModePerm,
	))
}
