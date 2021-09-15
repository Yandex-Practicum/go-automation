package snippetcomment

import "go/token"

type Comments struct {
	Comments    []Comment
	DocComments []DocComment
}

type Comment struct {
	StartPosition token.Pos
	Content       string
	Lines         []string
	IsDirective   bool
}

func NewComment(content string, lines []string, isDirective bool, startPosition token.Pos) Comment {
	return Comment{
		StartPosition: startPosition,
		Content:       content,
		Lines:         lines,
		IsDirective:   isDirective,
	}
}

type DocComment struct {
	StartPosition token.Pos
	Content       string
	Lines         []string
	EntitiesNames []string
	IsDirective   bool
}

func NewDocComment(content string, lines []string, entitiesNames []string, isDirective bool, pos token.Pos) DocComment {
	return DocComment{
		StartPosition: pos,
		Content:       content,
		Lines:         lines,
		EntitiesNames: entitiesNames,
		IsDirective:   isDirective,
	}
}
