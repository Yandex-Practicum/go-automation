package snippetsearch_test

import (
	"testing"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetsearch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindSnippets(t *testing.T) {
	t.Run("EmptyString", func(t *testing.T) {
		assert.Len(t, snippetsearch.FindSnippets(""), 0)
	})

	t.Run("NotClosedSnippet", func(t *testing.T) {
		assert.Len(t, snippetsearch.FindSnippets("```go и не закрыли сниппет"), 0)
	})

	t.Run("OneSnippet", func(t *testing.T) {
		snippets := snippetsearch.FindSnippets("```go\nа внутри go-код\n```")
		require.Len(t, snippets, 1)
		assert.EqualValues(t, snippetsearch.Snippet{
			Position: snippetsearch.SnippetPosition{
				Start: 6,
				End:   32,
			},
			Text: "а внутри go-код\n",
		}, snippets[0])
	})

	t.Run("MultipleSnippets", func(t *testing.T) {
		snippets := snippetsearch.FindSnippets("```go\nsnippet 1\n``` text ```go\nsnippet 2\n```")
		require.Len(t, snippets, 2)
		assert.EqualValues(t, "snippet 1\n", snippets[0].Text)
		assert.EqualValues(t, "snippet 2\n", snippets[1].Text)
	})
}
