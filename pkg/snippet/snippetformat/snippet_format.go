package snippetformat

import (
	"go/format"
	"strings"
)

func FormatSnippet(text string) (string, error) {
	b, err := format.Source([]byte(text))
	if err != nil {
		return "", err
	}

	// dirty hack to use 4 spaces in formatting instead of tabs
	// i hope no one uses plain tabs inside their strings in lessons
	return strings.ReplaceAll(string(b), "\t", "    "), nil
}
