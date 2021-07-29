package snippetformat

import (
	"fmt"
	"strings"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetsearch"
	"github.com/pkg/errors"
)

func FormatText(text string) (string, error) {
	snippets := snippetsearch.FindSnippets(text)
	if len(snippets) == 0 {
		return text, nil
	}

	var sb strings.Builder
	sb.Grow(len(text))

	sb.WriteString(getNonSnippetPrefix(text, snippets[0]))
	for i, snippet := range snippets {
		formattedSnippet, err := FormatSnippet(snippet.Text)
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("failed to format snippet \n%s\n:", snippet.Text))
		}
		sb.WriteString(formattedSnippet)

		if i != len(snippets)-1 {
			sb.WriteString(text[snippet.Position.End:snippets[i+1].Position.Start])
		}
	}

	sb.WriteString(getNonSnippetSuffix(text, snippets[len(snippets)-1]))

	return sb.String(), nil
}

func getNonSnippetPrefix(text string, firstSnippet snippetsearch.Snippet) string {
	if firstSnippet.Position.Start == 0 {
		return ""
	}

	return text[:firstSnippet.Position.Start]
}

func getNonSnippetSuffix(text string, lastSnippet snippetsearch.Snippet) string {
	if lastSnippet.Position.End == len(text) {
		return ""
	}

	return text[lastSnippet.Position.End:]
}
