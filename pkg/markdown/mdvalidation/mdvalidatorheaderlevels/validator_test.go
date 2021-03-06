package mdvalidatorheaderlevels_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/markdown/mdvalidation/mdvalidatorheaderlevels"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func TestValidator(t *testing.T) {
	v := mdvalidatorheaderlevels.NewValidator()

	testCases := []struct {
		Name                  string
		DocSource             string
		ExpectedErrorPrefixes []string
	}{
		{
			Name:      "First level header",
			DocSource: "# Header",
		},
		{
			Name:      "Second level header",
			DocSource: "## Header",
		},
		{
			Name:      "Third level header",
			DocSource: "### Header",
			ExpectedErrorPrefixes: []string{
				"Use only headers of level 1 and 2",
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
