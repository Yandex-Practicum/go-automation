package snippetformat_test

import (
	"testing"

	"github.com/Yandex-Practicum/go-automation/snippetformat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatText(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          string
		ExpectedOutput string
	}{
		{
			Name:           "NoSnippets",
			Input:          "text",
			ExpectedOutput: "text",
		},
		{
			Name:           "SingleSnippet",
			Input:          "prefix\n```go\nb.buf   =  b.buf[ :0  ]\n```\nsuffix",
			ExpectedOutput: "prefix\n```go\nb.buf = b.buf[:0]\n```\nsuffix",
		},
		{
			Name:           "MultipleSnippets",
			Input:          "prefix\n```go\nb.buf   =  b.buf[ :0  ]\n```\ninter snippet\n```go\nb.buf   =  b.buf[ :0  ]\n```\nsuffix",
			ExpectedOutput: "prefix\n```go\nb.buf = b.buf[:0]\n```\ninter snippet\n```go\nb.buf = b.buf[:0]\n```\nsuffix",
		},
		{
			Name:           "StartSnippet",
			Input:          "```go\nb.buf   =  b.buf[ :0  ]\n```\nsuffix",
			ExpectedOutput: "```go\nb.buf = b.buf[:0]\n```\nsuffix",
		},
		{
			Name:           "EndSnippet",
			Input:          "prefix\n```go\nb.buf   =  b.buf[ :0  ]\n```",
			ExpectedOutput: "prefix\n```go\nb.buf = b.buf[:0]\n```",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			output, err := snippetformat.FormatText(tc.Input)
			require.NoError(t, err)
			assert.EqualValues(t, tc.ExpectedOutput, output)
		})
	}
}
