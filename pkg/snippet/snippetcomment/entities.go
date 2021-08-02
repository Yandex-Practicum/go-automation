package snippetcomment

type Comments struct {
	Comments    []Comment
	DocComments []DocComment
}

type Comment struct {
	Content string
}

func NewComment(content string) Comment {
	return Comment{
		Content: content,
	}
}

type DocComment struct {
	Content string
}

func NewDocComment(content string) DocComment {
	return DocComment{
		Content: content,
	}
}
