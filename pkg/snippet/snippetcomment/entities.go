package snippetcomment

import "go/token"

type Comments struct {
	Comments    []Comment
	DocComments []DocComment
}

type Comment struct {
	StartPosition token.Pos
	Content       string
	IsDirective   bool
}

func NewComment(content string, isDirective bool, startPosition token.Pos) Comment {
	return Comment{
		StartPosition: startPosition,
		Content:       content,
		IsDirective:   isDirective,
	}
}

type DocComment struct {
	StartPosition token.Pos
	Content       string
	EntitiesNames []string
	IsDirective   bool
}

func NewDocComment(content string, entitiesNames []string, isDirective bool, pos token.Pos) DocComment {
	return DocComment{
		StartPosition: pos,
		Content:       content,
		EntitiesNames: entitiesNames,
		IsDirective:   isDirective,
	}
}
