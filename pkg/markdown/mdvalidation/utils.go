package mdvalidation

import (
	"strings"

	"github.com/yuin/goldmark/ast"
)

type NodeCallback func(node ast.Node) (err error, goDeeper bool)

func TraverseTree(root ast.Node, callback NodeCallback) error {
	if err, goDeeper := callback(root); err != nil || !goDeeper {
		return err
	}

	currChild := root.FirstChild()

	for i := 0; i != root.ChildCount(); i++ {
		if err := TraverseTree(currChild, callback); err != nil {
			return err
		}

		currChild = currChild.NextSibling()
	}

	return nil
}

func GetNodeText(node ast.Node, docSource []byte) string {
	if node.Type() == ast.TypeInline {
		return string(node.Text(docSource))
	}

	lines := node.Lines()

	var sb strings.Builder
	for _, segment := range lines.Sliced(0, lines.Len()) {
		sb.Write(segment.Value(docSource))
		if err := TraverseTree(node, func(node ast.Node) (err error, goDeeper bool) {
			sb.WriteString(GetNodeText(node, docSource))
			return nil, true
		}); err != nil {
			panic(err)
		}
	}

	return sb.String()
}
