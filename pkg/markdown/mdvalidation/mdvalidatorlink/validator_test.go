package mdvalidatorlink_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/markdown/mdvalidation/mdvalidatorlink"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func TestValidator(t *testing.T) {
	v := mdvalidatorlink.NewValidator()

	testCases := []struct {
		Name                  string
		DocSource             string
		ExpectedErrorPrefixes []string
	}{
		{
			Name:      "NoLink",
			DocSource: "text",
		},
		{
			Name:      "LinkWithTarget",
			DocSource: "[text](url){target=\"_blank\"}",
		},
		{
			Name:      "LinkWithoutTarget",
			DocSource: "[text](url)",
			ExpectedErrorPrefixes: []string{
				"All links must end up with {target=\"_blank\"}; Link content:",
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
