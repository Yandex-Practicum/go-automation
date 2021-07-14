package compilation

type Error struct {
	msg string
}

var _ error = (*Error)(nil)

func (e *Error) Error() string {
	return e.msg
}
