package snippetformat_test

import (
	"testing"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetformat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatSnippet(t *testing.T) {
	formatted, err := snippetformat.FormatSnippet("b.buf   =  b.buf[ :0  ]")
	require.NoError(t, err)
	assert.EqualValues(t, "b.buf = b.buf[:0]", formatted)
}
