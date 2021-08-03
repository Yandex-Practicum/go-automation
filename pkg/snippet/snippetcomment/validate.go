package snippetcomment

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func ValidateComment(comment Comment) error {
	content := comment.Content

	if len(content) == 0 {
		return errors.New("Empty comments are not allowed")
	}

	contentRunes := []rune(content)

	if lastRune := contentRunes[len(contentRunes)-1]; lastRune == '.' {
		return errors.New("Do not use . at the end of line comments")
	}

	if firstRune := contentRunes[0]; isRussianRune(firstRune) && isUpperCaseRune(firstRune) {
		return errors.New("First letter must be in lower case")
	}

	return nil
}

var (
	russianAlphabetLower = regexp.MustCompile("[абвгдеёжзийклмнопрстуфхцчшщъыьэюя]")
	russianAlphabetUpper = regexp.MustCompile("[АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ]")
)

func isRussianRune(r rune) bool {
	return russianAlphabetUpper.MatchString(string(r)) || russianAlphabetLower.MatchString(string(r))
}

func isUpperCaseRune(r rune) bool {
	runeString := string(r)
	return strings.ToLower(runeString) != runeString
}
