package snippetcomment

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetparse"
)

func ExtractComments(snippet snippetparse.ParsedSnippet) (Comments, error) {
	docCommentWithPosition := extractDocComments(snippet)

	comments := extractSimpleComments(snippet, docCommentWithPosition)

	return Comments{
		DocComments: collectDocComments(docCommentWithPosition),
		Comments:    comments,
	}, nil
}

func extractSimpleComments(snippet snippetparse.ParsedSnippet, docComments []docCommentWithPosition) []Comment {
	knownDocCommentPositions := make(map[token.Pos]struct{}, len(docComments))
	for _, c := range docComments {
		knownDocCommentPositions[c.StartPosition] = struct{}{}
	}

	var result []Comment
	for _, comment := range snippet.AST.Comments {
		if _, ok := knownDocCommentPositions[comment.Pos()]; ok {
			continue
		}

		result = append(result, newNormalizedComment(comment.Text()))
	}

	return result
}

func newNormalizedComment(content string) Comment {
	return NewComment(normalizeComment(content))
}

type docCommentWithPosition struct {
	Comment       DocComment
	StartPosition token.Pos
}

func collectDocComments(comments []docCommentWithPosition) []DocComment {
	var result []DocComment
	for _, comment := range comments {
		result = append(result, comment.Comment)
	}

	return result
}

func extractDocComments(snippet snippetparse.ParsedSnippet) []docCommentWithPosition {
	var result []docCommentWithPosition

	if packageDoc := snippet.AST.Doc; packageDoc != nil {
		result = append(result, newNormalizedDocCommentWithPosition(packageDoc.Text(), packageDoc.Pos()))
	}

	for _, decl := range snippet.AST.Decls {
		switch typedDecl := decl.(type) {
		case *ast.GenDecl:
			if typedDecl.Doc != nil {
				result = append(result, newNormalizedDocCommentWithPosition(typedDecl.Doc.Text(), typedDecl.Doc.Pos()))
			}

		case *ast.FuncDecl:
			if typedDecl.Doc != nil {
				result = append(result, newNormalizedDocCommentWithPosition(typedDecl.Doc.Text(), typedDecl.Doc.Pos()))
			}
		}
	}

	return result
}

func newNormalizedDocCommentWithPosition(content string, pos token.Pos) docCommentWithPosition {
	return docCommentWithPosition{
		Comment:       NewDocComment(normalizeComment(content)),
		StartPosition: pos,
	}
}

func normalizeComment(content string) string {
	return strings.TrimSpace(content)
}
