package namevalidation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Yandex-Practicum/go-automation/filesearch"
	"github.com/pkg/errors"
)

func Validate(rootDir string) error {
	var errMsgs []string

	invalidDirNames, err := getInvalidLessonNames(rootDir)
	if err != nil {
		return err
	}
	if len(invalidDirNames) > 0 {
		errMsgs = append(errMsgs,
			fmt.Sprintf(
				"Lesson dirs %v do not match lesson dir naming convention. Lesson dir name must consist of two digits (for example '01')", invalidDirNames,
			),
		)
	}

	invalidFileNames, err := getInvalidMDFileNames(rootDir)
	if err != nil {
		return err
	}
	if len(invalidFileNames) > 0 {
		errMsgs = append(errMsgs,
			fmt.Sprintf(
				"Files %v do not match markdown naming convention. Lesson file must have format like '01-Lorem-ipsum' (two digits, first word starts from upper case letter)", invalidFileNames,
			),
		)
	}

	if len(errMsgs) > 0 {
		return errors.New(strings.Join(errMsgs, "\n"))
	}

	return nil
}

func getInvalidMDFileNames(rootDir string) ([]string, error) {
	mdFileNames, err := filesearch.GetFileWithExtensionPaths(rootDir, "md")
	if err != nil {
		return nil, err
	}

	var invalidFileNames []string
	for _, fileName := range mdFileNames {
		if !isValidFileName(getLastPathPart(fileName)) {
			invalidFileNames = append(invalidFileNames, fileName)
		}
	}

	return invalidFileNames, nil
}

func getInvalidLessonNames(rootDir string) ([]string, error) {
	lessonDirNames, err := getLessonDirNames(rootDir, 1)
	if err != nil {
		return nil, err
	}

	var invalidDirNames []string
	for _, dirName := range lessonDirNames {
		if !isValidLessonDirName(getLastPathPart(dirName)) {
			invalidDirNames = append(invalidDirNames, dirName)
		}
	}

	return invalidDirNames, nil
}

func getLessonDirNames(rootDir string, levelNum int) ([]string, error) {
	var paths []string
	if err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		pathSuffix := strings.TrimPrefix(path, rootDir)
		if pathSuffix == "" {
			return nil
		}

		pathSuffix = strings.TrimPrefix(pathSuffix, "/")

		suffixParts := strings.Split(pathSuffix, "/")
		if len(suffixParts) == 0 || len(suffixParts) > levelNum {
			return nil
		}

		paths = append(paths, path)
		return nil
	}); err != nil {
		return nil, err
	}

	return paths, nil
}

var lessonDirNameRegexp = regexp.MustCompile(`^[0-9]{2}$`)

func isValidLessonDirName(dirName string) bool {
	return lessonDirNameRegexp.MatchString(dirName)
}

var filenameRegexp = regexp.MustCompile(`^[0-9]{2}-[A-Z][a-z]*(-[a-z]+)*\.md$`)

func isValidFileName(fileName string) bool {
	return filenameRegexp.MatchString(fileName)
}

func getLastPathPart(path string) string {
	pathParts := strings.Split(path, "/")
	return pathParts[len(pathParts)-1]
}
