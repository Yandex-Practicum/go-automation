package snippetcomment

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetparse"
)

func ExtractComments(snippet snippetparse.ParsedSnippet) Comments {
	docComments := extractDocComments(snippet)

	comments := extractSimpleComments(snippet, docComments)

	return Comments{
		DocComments: docComments,
		Comments:    comments,
	}
}

func extractSimpleComments(snippet snippetparse.ParsedSnippet, docComments []DocComment) []Comment {
	knownDocCommentPositions := make(map[token.Pos]struct{}, len(docComments))
	for _, c := range docComments {
		knownDocCommentPositions[c.StartPosition] = struct{}{}
	}

	var result []Comment
	for _, comment := range snippet.AST.Comments {
		if _, ok := knownDocCommentPositions[comment.Pos()]; ok {
			continue
		}

		commentText, isDirective := getCommentText(comment)
		result = append(result, newNormalizedComment(commentText, isDirective, comment.Pos()))
	}

	return result
}

func newNormalizedComment(content string, isDirective bool, pos token.Pos) Comment {
	return NewComment(normalizeComment(content), isDirective, pos)
}

func extractDocComments(snippet snippetparse.ParsedSnippet) []DocComment {
	var result []DocComment

	handleDocCommentGroup := func(g *ast.CommentGroup, decl ast.Decl) {
		commentText, isDirective := getCommentText(g)
		result = append(result, newNormalizedDocComment(commentText, extractDeclarationNames(decl), isDirective, g.Pos()))
	}

	if packageDoc := snippet.AST.Doc; packageDoc != nil {
		commentText, isDirective := getCommentText(packageDoc)
		result = append(result, newNormalizedDocComment(commentText, namesFromIdents(snippet.AST.Name), isDirective, packageDoc.Pos()))
	}

	for _, decl := range snippet.AST.Decls {
		switch typedDecl := decl.(type) {
		case *ast.GenDecl:
			if typedDecl.Doc != nil {
				handleDocCommentGroup(typedDecl.Doc, decl)
			}

		case *ast.FuncDecl:
			if typedDecl.Doc != nil {
				handleDocCommentGroup(typedDecl.Doc, decl)
			}
		}
	}

	return result
}

func getCommentText(g *ast.CommentGroup) (string, bool) {
	commentText := g.Text()
	if len(commentText) > 0 {
		return commentText, false
	}

	for _, line := range g.List {
		if isDirectiveComment(line.Text) {
			return "", true
		}
	}

	return "", false
}

func isDirectiveComment(c string) bool {
	switch c[1] {
	case '/':
		//-style comment (no newline at the end)
		c = c[2:]
		if len(c) == 0 {
			// empty line
			return false
		}
		if c[0] == ' ' {
			// strip first space - required for Example tests
			c = c[1:]
			return false
		}

		return isDirective(c)
	default:
		return false
	}
}

// isDirective reports whether c is a comment directive.
func isDirective(c string) bool {
	// "//line " is a line directive.
	// (The // has been removed.)
	if strings.HasPrefix(c, "line ") {
		return true
	}

	// "//[a-z0-9]+:[a-z0-9]"
	// (The // has been removed.)
	colon := strings.Index(c, ":")
	if colon <= 0 || colon+1 >= len(c) {
		return false
	}
	for i := 0; i <= colon+1; i++ {
		if i == colon {
			continue
		}
		b := c[i]
		if !('a' <= b && b <= 'z' || '0' <= b && b <= '9') {
			return false
		}
	}
	return true
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

func newNormalizedDocComment(content string, entitiesNames []string, isDirective bool, pos token.Pos) DocComment {
	return NewDocComment(normalizeComment(content), entitiesNames, isDirective, pos)
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
