package snippetparse

import (
	"go/ast"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetsearch"
)

type ParsedSnippet struct {
	Snippet snippetsearch.Snippet
	AST     *ast.File
}
