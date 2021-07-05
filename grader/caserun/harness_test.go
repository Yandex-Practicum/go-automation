package caserun_test

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	moduleDir         string
	originalModuleDir string
}

func (s *Suite) SetupTest() {
	var err error

	s.moduleDir, err = ioutil.TempDir("", "caserun_*")
	s.Require().NoError(err)

	s.originalModuleDir, err = ioutil.TempDir("", "caserun_original_*")
	s.Require().NoError(err)
}

func (s *Suite) TearDownTest() {
	s.Require().NoError(os.RemoveAll(s.moduleDir))
}

func (s *Suite) CreateMain(content, dir string) {
	s.CreateFile("main.go", content, dir)
}

func (s *Suite) CreateMod(dir string) {
	s.CreateFile("go.mod", `
module exercise

go 1.15
`, dir)
}

func (s *Suite) CreateFile(name, content, dir string) {
	s.Require().NoError(ioutil.WriteFile(
		path.Join(dir, name),
		[]byte(content),
		os.ModePerm,
	))
}
