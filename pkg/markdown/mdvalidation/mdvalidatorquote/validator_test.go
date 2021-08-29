package mdvalidatorquote_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/markdown/mdvalidation/mdvalidatorquote"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func TestValidator(t *testing.T) {
	v := mdvalidatorquote.NewValidator()

	testCases := []struct {
		Name                  string
		DocSource             string
		ExpectedErrorPrefixes []string
	}{
		{
			Name:      "NoQuote",
			DocSource: "# Header",
		},
		{
			Name:      "Quote",
			DocSource: "> quote",
			ExpectedErrorPrefixes: []string{
				"Block quote is banned; https://yandex-edu.slack.com/archives/C020W0HAH2Q/p1628930940064700; Quote content:",
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
