package mdvalidatoritalic_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/markdown/mdvalidation/mdvalidatoritalic"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func TestValidator(t *testing.T) {
	v := mdvalidatoritalic.NewValidator()

	testCases := []struct {
		Name                  string
		DocSource             string
		ExpectedErrorPrefixes []string
	}{
		{
			Name:      "NoStars",
			DocSource: "text",
		},
		{
			Name:      "SpaceInTheBeginning",
			DocSource: "* text*",
		},
		{
			Name:      "SpaceInTheEnd",
			DocSource: "*text *",
		},
		{
			Name:      "Bold",
			DocSource: "**text**",
		},
		{
			Name:      "CodeSpan",
			DocSource: "`*text*`",
		},
		{
			Name:      "CodeBlock",
			DocSource: "```\n*text*\n```",
		},
		{
			Name:      "Italic",
			DocSource: "*text*",
			ExpectedErrorPrefixes: []string{
				"Italic is banned use eather backticks or bold; Node content:",
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
