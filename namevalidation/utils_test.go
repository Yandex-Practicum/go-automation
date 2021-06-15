package namevalidation_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/Yandex-Practicum/go-automation/namevalidation"
	"github.com/stretchr/testify/suite"
)

type ValidationTestSuite struct {
	suite.Suite

	rootDir string
}

func TestValidationTestSuite(t *testing.T) {
	suite.Run(t, &ValidationTestSuite{})
}

func (s *ValidationTestSuite) SetupTest() {
	rootDir, err := ioutil.TempDir(os.TempDir(), "name_validation_test_*")
	s.Require().NoError(err)

	s.rootDir = rootDir
}

func (s *ValidationTestSuite) TestBadLessonDir() {
	s.createLessonDir("bad_dir")
	err := namevalidation.Validate(s.rootDir)
	s.Require().Error(err)

	s.EqualValues(fmt.Sprintf("Lesson dirs [%s/bad_dir] do not match lesson dir naming convention. Lesson dir name must consist of two digits (for example '01')", s.rootDir), err.Error())
}

func (s *ValidationTestSuite) TestOkLessonDir() {
	s.createLessonDir("02")
	s.Require().NoError(namevalidation.Validate(s.rootDir))
}

func (s *ValidationTestSuite) TestBadFileName() {
	s.createLessonFile("bad_lesson_file.md")
	err := namevalidation.Validate(s.rootDir)
	s.Require().Error(err)

	s.EqualValues(fmt.Sprintf("Files [%s/bad_lesson_file.md] do not match markdown naming convention. Lesson file must have format like '01-Lorem-ipsum' (two digits, first word starts from upper case letter)", s.rootDir), err.Error())
}

func (s *ValidationTestSuite) TestOkFileName() {
	s.createLessonFile("01-Ok-file.md")
	s.Require().NoError(namevalidation.Validate(s.rootDir))
}

func (s *ValidationTestSuite) createLessonDir(name string) {
	s.Require().NoError(os.Mkdir(path.Join(s.rootDir, name), os.ModePerm))
}

func (s *ValidationTestSuite) createLessonFile(name string) {
	_, err := os.Create(path.Join(s.rootDir, name))
	s.Require().NoError(err)
}
