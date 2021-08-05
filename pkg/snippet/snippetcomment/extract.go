package snippetcomment

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetparse"
)

func ExtractComments(snippet snippetparse.ParsedSnippet) Comments {
	docCommentWithPosition := extractDocComments(snippet)

	comments := extractSimpleComments(snippet, docCommentWithPosition)

	return Comments{
		DocComments: collectDocComments(docCommentWithPosition),
		Comments:    comments,
	}
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
		result = append(result, newNormalizedDocCommentWithPosition(packageDoc.Text(), namesFromIdents(snippet.AST.Name), packageDoc.Pos()))
	}

	for _, decl := range snippet.AST.Decls {
		switch typedDecl := decl.(type) {
		case *ast.GenDecl:
			if typedDecl.Doc != nil {
				result = append(result, newNormalizedDocCommentWithPosition(typedDecl.Doc.Text(), extractDeclarationNames(decl), typedDecl.Doc.Pos()))
			}

		case *ast.FuncDecl:
			if typedDecl.Doc != nil {
				result = append(result, newNormalizedDocCommentWithPosition(typedDecl.Doc.Text(), extractDeclarationNames(decl), typedDecl.Doc.Pos()))
			}
		}
	}

	return result
}

func extractDeclarationNames(decl ast.Decl) []string {
	switch typedDecl := decl.(type) {
	case *ast.FuncDecl:
		return []string{typedDecl.Name.Name}

	case *ast.GenDecl:
		specs := typedDecl.Specs

		var result []string
		for _, spec := range specs {
			switch typedSpec := spec.(type) {
			case *ast.TypeSpec:
				result = append(result, namesFromIdents(typedSpec.Name)...)

			case *ast.ValueSpec:
				result = append(result, namesFromIdents(typedSpec.Names...)...)
			}
		}

		return result

	default:
		return nil
	}
}

func newNormalizedDocCommentWithPosition(content string, entitiesNames []string, pos token.Pos) docCommentWithPosition {
	return docCommentWithPosition{
		Comment:       NewDocComment(normalizeComment(content), entitiesNames),
		StartPosition: pos,
	}
}

func normalizeComment(content string) string {
	return strings.TrimSpace(content)
}

func namesFromIdents(idents ...*ast.Ident) []string {
	result := make([]string, 0, len(idents))
	for _, ident := range idents {
		if ident != nil {
			result = append(result, ident.Name)
		}
	}

	return result
}
