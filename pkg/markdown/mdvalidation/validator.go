package mdvalidation

import (
	"context"

	"github.com/yuin/goldmark/ast"
)

type Validator interface {
	ValidateTree(ctx context.Context, tree ast.Node, docSource []byte) ([]*ValidationInfo, error)
}

type ValidationInfo struct {
	Message string
	Node    ast.Node
}
