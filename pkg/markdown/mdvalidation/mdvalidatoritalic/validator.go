package mdvalidatoritalic

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/markdown/mdvalidation"
	"github.com/yuin/goldmark/ast"
)

func NewValidator() mdvalidation.Validator {
	return &validator{}
}

type validator struct{}

var _ mdvalidation.Validator = (*validator)(nil)

func (v *validator) ValidateTree(ctx context.Context, tree ast.Node, docSource []byte) ([]*mdvalidation.ValidationInfo, error) {
	badNodes, err := v.getBadNodes(ctx, tree, docSource)
	if err != nil || len(badNodes) == 0 {
		return nil, err
	}

	validationResults := make([]*mdvalidation.ValidationInfo, 0, len(badNodes))
	for _, node := range badNodes {
		validationResults = append(validationResults, &mdvalidation.ValidationInfo{
			Message: fmt.Sprintf("Italic is banned use eather backticks or bold; Node content: \n%s",
				mdvalidation.GetNodeText(node, docSource, mdvalidation.NodeTextOptionLimit(200))),
			Node: node,
		})
	}

	return validationResults, nil
}

func (v *validator) getBadNodes(ctx context.Context, tree ast.Node, docSource []byte) ([]ast.Node, error) {
	var result []ast.Node
	if err := mdvalidation.TraverseTree(tree, func(node ast.Node) (error, bool) {
		switch node.(type) {
		case *ast.CodeSpan, *ast.CodeBlock:
			return nil, false

		case *ast.Document:
			return nil, true

		default:
			if v.nodeContainsItalic(ctx, mdvalidation.GetNodeText(node, docSource)) {
				result = append(result, node)
			}
		}

		return nil, false
	}); err != nil {
		return nil, err
	}

	return result, nil
}

var mayBeItalicRegexp = regexp.MustCompile(`(?P<Snippet>\*[^*\s].*[^*\s]\*)`)

func (v *validator) nodeContainsItalic(ctx context.Context, nodeText string) bool {
	submatches := mayBeItalicRegexp.FindAllStringSubmatchIndex(nodeText, -1)
	for _, submatch := range submatches {
		start, end := submatch[2], submatch[3]

		isBoldStart := start != 0 && nodeText[start-1] == '*'
		isBoldEnd := end != len(nodeText) && nodeText[end] == '*'

		if isBoldStart && isBoldEnd {
			return false
		}

		return true
	}

	return false
}
