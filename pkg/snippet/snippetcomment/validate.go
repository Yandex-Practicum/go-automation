package snippetcomment

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func ValidateDocComment(comment DocComment) error {
	if comment.IsDirective {
		return nil
	}

	content := comment.Content

	if len(content) == 0 {
		return newEmptyCommentError()
	}

	docPrefix := getCommonDocNamesPrefix(comment.EntitiesNames)
	if !strings.HasPrefix(content, docPrefix) {
		return errors.New(fmt.Sprintf("Doc comment must start with documented entity name (need prefix %s)", docPrefix))
	}

	contentRunes := []rune(content)
	if lastRune := contentRunes[len(contentRunes)-1]; lastRune != '.' {
		return errors.New("Doc comments must end up with .")
	}

	if hasTooLingLine(comment.Lines) {
		return errors.New(fmt.Sprintf("Every line of comment must be shorter than %d symbols", maxCommentLineLength))
	}

	return nil
}

func getCommonDocNamesPrefix(names []string) string {
	if len(names) == 0 {
		return ""
	}

	nameRunes := make([][]rune, 0, len(names))
	for _, name := range names {
		nameRunes = append(nameRunes, []rune(name))
	}

	minRuneLength := getMinLength(nameRunes)

	var commonRunePrefix []rune
	for i := 0; i != minRuneLength; i++ {
		if firstNameRune := nameRunes[0][i]; haveSameRuneOnPosition(nameRunes, i) {
			commonRunePrefix = append(commonRunePrefix, firstNameRune)
		}
	}

	return string(commonRunePrefix)
}

func haveSameRuneOnPosition(runeArrays [][]rune, pos int) bool {
	firstArrayRune := runeArrays[0][pos]
	for _, arr := range runeArrays {
		if arr[pos] != firstArrayRune {
			return false
		}
	}

	return true
}

func getMinLength(runeArrays [][]rune) int {
	result := len(runeArrays[0])
	for _, arr := range runeArrays {
		if len(arr) < result {
			result = len(arr)
		}
	}

	return result
}

func ValidateComment(comment Comment) error {
	if comment.IsDirective {
		return nil
	}

	content := comment.Content

	if len(content) == 0 {
		return newEmptyCommentError()
	}

	if strings.HasSuffix(content, ".") && !strings.HasSuffix(content, "...") {
		return errors.New("Do not use . at the end of line comments")
	}

	firstRune := []rune(content)[0]

	if isUpperCaseRussianRune(firstRune) {
		return errors.New("First letter must be in lower case")
	}

	if hasTooLingLine(comment.Lines) {
		return errors.New(fmt.Sprintf("Every line of comment must be shorter than %d symbols", maxCommentLineLength))
	}

	return nil
}

func hasTooLingLine(lines []string) bool {
	for _, line := range lines {
		if isTooLongCommentLine(line) {
			return true
		}
	}

	return false
}

const maxCommentLineLength = 85

func isTooLongCommentLine(commentLine string) bool {
	commentLineRunes := []rune(commentLine)
	return len(commentLineRunes) > maxCommentLineLength
}

func newEmptyCommentError() error {
	return errors.New("Empty comments are not allowed")
}

var (
	russianAlphabetLower = regexp.MustCompile("[абвгдеёжзийклмнопрстуфхцчшщъыьэюя]")
)

func isRussianRune(r rune) bool {
	return russianAlphabetLower.MatchString(strings.ToLower(string(r)))
}

func isUpperCaseRussianRune(r rune) bool {
	runeString := string(r)
	return isRussianRune(r) && strings.ToLower(string(r)) != runeString
}
