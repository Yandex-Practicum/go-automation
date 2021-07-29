package snippetformat

import "go/format"

func FormatSnippet(text string) (string, error) {
	b, err := format.Source([]byte(text))
	if err != nil {
		return "", err
	}

	return string(b), nil
}
