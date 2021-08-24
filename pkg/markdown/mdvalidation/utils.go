package mdvalidation

import (
	"math"

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

	nodeRange := GetNodeRange(node)

	return string(docSource[nodeRange.Start:nodeRange.Stop])
}

type NodeRange struct {
	Start int
	Stop  int
}

func GetNodeRange(node ast.Node) NodeRange {
	if node.Type() == ast.TypeInline {
		panic("inline nodes have no range")
	}

	start := math.MaxInt32
	stop := -1
	if err := TraverseTree(node, func(node ast.Node) (err error, goDeeper bool) {
		if node.Type() == ast.TypeInline {
			return nil, true
		}

		lines := node.Lines()
		for _, line := range lines.Sliced(0, lines.Len()) {
			if line.Start < start {
				start = line.Start
			}

			if line.Stop > stop {
				stop = line.Stop
			}
		}

		return nil, true
	}); err != nil {
		panic(err)
	}

	return NodeRange{
		Start: start,
		Stop:  stop,
	}
}
