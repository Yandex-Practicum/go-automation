package mdvalidatorlink

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
	badLinks, err := v.getBadLinks(ctx, tree, docSource)
	if err != nil || len(badLinks) == 0 {
		return nil, err
	}

	validationResults := make([]*mdvalidation.ValidationInfo, 0, len(badLinks))
	for _, node := range badLinks {
		validationResults = append(validationResults, &mdvalidation.ValidationInfo{
			Message: fmt.Sprintf("All links must end up with {target=\"_blank\"}; Link content: \n%s",
				mdvalidation.GetNodeText(node, docSource, mdvalidation.NodeTextOptionLimit(200))),
			Node: node,
		})
	}

	return validationResults, nil
}

var targetPrefixRegexp = regexp.MustCompile(`^\(.*\)\{target="_blank"\}.*`)

func (v *validator) getBadLinks(ctx context.Context, tree ast.Node, docSource []byte) ([]ast.Node, error) {
	var result []ast.Node
	if err := mdvalidation.TraverseTree(tree, func(node ast.Node) (error, bool) {
		switch node.(type) {
		case *ast.Link:
			nodeStop := v.getLinkStopPosition(node.(*ast.Link))
			if nodeStop == len(docSource) || !targetPrefixRegexp.Match(docSource[nodeStop:]) {
				result = append(result, node)
			}

			return nil, false
		}

		return nil, true
	}); err != nil {
		return nil, err
	}

	return result, nil
}

// it's a trash, i know. But parsing lib is a little imperfect
func (v *validator) getLinkStopPosition(node *ast.Link) int {
	var textNode *ast.Text
	textNodeCandidate := node.FirstChild()
	for i := 0; i != 10; i++ {
		if _, ok := textNodeCandidate.(*ast.Text); ok {
			textNode = textNodeCandidate.(*ast.Text)
			break
		} else {
			textNodeCandidate = textNodeCandidate.FirstChild()
		}
	}

	if textNode == nil {
		panic("link without text node")
	}

	return textNode.Segment.Stop + 1
}
