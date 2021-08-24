package mdvalidatorheaderlevels

import (
	"context"
	"fmt"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/markdown/mdvalidation"
	"github.com/yuin/goldmark/ast"
)

func NewValidator() mdvalidation.Validator {
	return &validator{}
}

type validator struct{}

var _ mdvalidation.Validator = (*validator)(nil)

func (v *validator) ValidateTree(ctx context.Context, tree ast.Node, docSource []byte) ([]*mdvalidation.ValidationInfo, error) {
	badHeaders, err := v.getBadHeaders(ctx, tree)
	if err != nil || len(badHeaders) == 0 {
		return nil, err
	}

	validationResults := make([]*mdvalidation.ValidationInfo, 0, len(badHeaders))
	for _, node := range badHeaders {
		validationResults = append(validationResults, &mdvalidation.ValidationInfo{
			Message: fmt.Sprintf("Use only headers of level 1 and 2; Header content: \n%s", mdvalidation.GetNodeText(node, docSource)),
			Node:    node,
		})
	}

	return validationResults, nil
}

func (v *validator) getBadHeaders(ctx context.Context, tree ast.Node) ([]ast.Node, error) {
	var result []ast.Node
	if err := mdvalidation.TraverseTree(tree, func(node ast.Node) (error, bool) {
		switch typedNode := node.(type) {
		case *ast.Heading:
			if typedNode.Level > 2 {
				result = append(result, node)
			}
		}

		return nil, true
	}); err != nil {
		return nil, err
	}

	return result, nil
}
