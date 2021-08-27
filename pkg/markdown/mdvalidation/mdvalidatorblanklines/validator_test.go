package mdvalidatorblanklines_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/markdown/mdvalidation/mdvalidatorblanklines"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func TestValidator(t *testing.T) {
	v := mdvalidatorblanklines.NewValidator()

	testCases := []struct {
		Name                  string
		DocSource             string
		ExpectedErrorPrefixes []string
	}{
		{
			Name:      "CodeBlockHasNoLines",
			DocSource: "# Header\n```go\nfoo\n```",
			ExpectedErrorPrefixes: []string{
				"Nodes FencedCodeBlock",
			},
		},
		{
			Name:      "HeadingHasNoLines",
			DocSource: "# Header\n# Header",
			ExpectedErrorPrefixes: []string{
				"Nodes Heading",
			},
		},
		{
			Name:      "QuoteHasNoLines",
			DocSource: "# Header\n> quote\ncontinued",
			ExpectedErrorPrefixes: []string{
				"Nodes Blockquote",
			},
		},
		{
			Name:      "ListHasNoLines",
			DocSource: "# Header\n- list item",
			ExpectedErrorPrefixes: []string{
				"Nodes List",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			p := goldmark.DefaultParser()

			docSource := []byte(tc.DocSource)
			infos, err := v.ValidateTree(context.Background(), p.Parse(text.NewReader(docSource)), docSource)
			require.NoError(t, err)

			require.EqualValues(t, len(tc.ExpectedErrorPrefixes), len(infos))

			for i, prefix := range tc.ExpectedErrorPrefixes {
				require.True(t, strings.HasPrefix(infos[i].Message, prefix), "prefix: %s\nmessage: %s", prefix, infos[i].Message)
			}
		})
	}
}
