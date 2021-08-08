package snippetparse_test

import (
	"strings"
	"testing"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetparse"
	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetsearch"
	"github.com/stretchr/testify/require"
)

func TestParseSnippet(t *testing.T) {
	snippetText := `
b.buf = b.buf[:] // this is assignment
b.buf = append(b.buf, 1)
`
	snippet, err := snippetparse.ParseSnippet(snippetsearch.Snippet{
		Text: snippetText,
	})
	require.NoError(t, err)

	comments := snippet.AST.Comments
	require.Len(t, comments, 1)
	require.EqualValues(t, "this is assignment", strings.TrimSpace(comments[0].Text()))
}
