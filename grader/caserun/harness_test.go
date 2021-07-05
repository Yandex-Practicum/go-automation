package caserun_test

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	tempDir string
}

func (s *Suite) SetupTest() {
	var err error
	s.tempDir, err = ioutil.TempDir("", "caserun_*")
	s.Require().NoError(err)
}

func (s *Suite) TearDownTest() {
	s.Require().NoError(os.RemoveAll(s.tempDir))
}

func (s *Suite) CreateMain(content string) {
	s.CreateFile("main.go", content)
}

func (s *Suite) CreateMod() {
	s.CreateFile("go.mod", `
module exercise

go 1.15
`)
}

func (s *Suite) CreateFile(name, content string) {
	s.Require().NoError(ioutil.WriteFile(
		path.Join(s.tempDir, name),
		[]byte(content),
		os.ModePerm,
	))
}
