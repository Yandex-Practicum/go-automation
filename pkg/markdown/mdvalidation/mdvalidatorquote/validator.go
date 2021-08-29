package mdvalidatorquote

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
	citationNodes, err := v.getCitationNodes(ctx, tree)
	if err != nil || len(citationNodes) == 0 {
		return nil, err
	}

	validationResults := make([]*mdvalidation.ValidationInfo, 0, len(citationNodes))
	for _, node := range citationNodes {
		validationResults = append(validationResults, &mdvalidation.ValidationInfo{
			Message: fmt.Sprintf("Block quote is banned; https://yandex-edu.slack.com/archives/C020W0HAH2Q/p1628930940064700; Quote content: \n%s",
				mdvalidation.GetNodeText(node, docSource, mdvalidation.NodeTextOptionLimit(200))),
			Node: node,
		})
	}

	return validationResults, nil
}

func (v *validator) getCitationNodes(ctx context.Context, tree ast.Node) ([]ast.Node, error) {
	var result []ast.Node
	if err := mdvalidation.TraverseTree(tree, func(node ast.Node) (error, bool) {
		switch node.(type) {
		case *ast.Blockquote:
			result = append(result, node)
		}

		return nil, true
	}); err != nil {
		return nil, err
	}

	return result, nil
}
