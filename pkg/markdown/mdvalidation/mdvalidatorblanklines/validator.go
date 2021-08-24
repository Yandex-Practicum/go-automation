package mdvalidatorblanklines

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
	badNodes, err := v.getBadNodes(ctx, tree)
	if err != nil || len(badNodes) == 0 {
		return nil, err
	}

	validationResults := make([]*mdvalidation.ValidationInfo, 0, len(badNodes))
	for _, node := range badNodes {
		validationResults = append(validationResults, &mdvalidation.ValidationInfo{
			Message: fmt.Sprintf("Nodes %s are expected to have blank lines above. Node content\n%s", node.Kind(), mdvalidation.GetNodeText(node, docSource)),
			Node:    node,
		})
	}

	return validationResults, nil
}

func (v *validator) getBadNodes(ctx context.Context, tree ast.Node) ([]ast.Node, error) {
	var result []ast.Node
	isFirstBlockNode := true
	if err := mdvalidation.TraverseTree(tree, func(node ast.Node) (error, bool) {
		if !v.isBlockNode(node) {
			return nil, true
		}

		if v.isBadNode(node, isFirstBlockNode) {
			result = append(result, node)
		}
		isFirstBlockNode = false

		return nil, false
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func (v *validator) isBadNode(node ast.Node, isFirstNode bool) bool {
	if isFirstNode {
		return false
	}

	return !node.HasBlankPreviousLines()
}

func (v *validator) isBlockNode(node ast.Node) bool {
	return node.Type() == ast.TypeBlock
}
