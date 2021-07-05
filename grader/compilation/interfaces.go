package compilation

import (
	"context"
)

type Compiler interface {
	CompilePackage(ctx context.Context, query Query) error
}
