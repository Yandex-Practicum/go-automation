package snippetcomment

import "go/token"

type Comments struct {
	Comments    []Comment
	DocComments []DocComment
}

type Comment struct {
	StartPosition token.Pos
	Content       string
}

func NewComment(content string, startPosition token.Pos) Comment {
	return Comment{
		StartPosition: startPosition,
		Content:       content,
	}
}

type DocComment struct {
	StartPosition token.Pos
	Content       string
	EntitiesNames []string
}

func NewDocComment(content string, entitiesNames []string, pos token.Pos) DocComment {
	return DocComment{
		StartPosition: pos,
		Content:       content,
		EntitiesNames: entitiesNames,
	}
}
