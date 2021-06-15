package wordvalidation_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/Yandex-Practicum/go-automation/wordvalidation"
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

func (s *ValidationTestSuite) TestWithStopWords() {
	s.createLesson("по-моему, тут есть стоп-слова")
	err := wordvalidation.Validate(s.rootDir)
	s.Require().Error(err)

	s.EqualValues(fmt.Sprintf("stop words found:\nfile %s has stop words [по-моему]\nhttps://smsreader.github.io/stopwords/ for convenient UI", s.getLessonFileName()), err.Error())
}

func (s *ValidationTestSuite) TestNoStopWords() {
	s.createLesson("а тут их нет")
	s.Require().NoError(wordvalidation.Validate(s.rootDir))
}

func (s *ValidationTestSuite) createLesson(content string) {
	s.Require().NoError(ioutil.WriteFile(s.getLessonFileName(), []byte(content), os.ModePerm))
}

func (s *ValidationTestSuite) getLessonFileName() string {
	return path.Join(s.rootDir, "file.md")
}
