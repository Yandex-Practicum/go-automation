package commandrun

import "context"

type Runner interface {
	Run(ctx context.Context) (*RunResult, error)
}
