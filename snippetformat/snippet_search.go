package snippetformat

import (
	"regexp"
)

type Snippet struct {
	Position SnippetPosition
	Text     string
}

type SnippetPosition struct {
	Start int
	End   int
}

var snippetRegexp = regexp.MustCompile("\\s*```[gG]o\\n(?P<Snippet>[\\s\\S]*?)```")

func FindSnippets(text string) []Snippet {
	submatches := snippetRegexp.FindAllStringSubmatchIndex(text, -1)

	result := make([]Snippet, 0, len(submatches))
	for _, submatch := range submatches {
		start, end := submatch[2], submatch[3]
		result = append(result, Snippet{
			Position: SnippetPosition{
				Start: start,
				End:   end,
			},
			Text: text[start:end],
		})
	}

	return result
}
