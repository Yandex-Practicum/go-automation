package compilation

import (
	"context"
	"time"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/grader/commandrun"
)

type Compiler interface {
	CompilePackage(ctx context.Context, query Query) error
}

type compiler struct {
}

var _ Compiler = (*compiler)(nil)

func NewCompiler() *compiler {
	return &compiler{}
}

func (c *compiler) CompilePackage(ctx context.Context, query Query) error {
	runner := c.getRunner(query)

	runInfo, err := runner.Run(ctx, commandrun.RunOptions{
		Dir: query.ModulePath,
	})
	if err != nil {
		return err
	}

	if runInfo.ExitCode != 0 {
		return &Error{
			msg: runInfo.Stderr,
		}
	}

	return nil
}

func (c *compiler) getRunner(query Query) commandrun.Runner {
	if query.BinaryPath == "" {
		return commandrun.NewRunner(time.Minute, "go", "build", ".")
	}

	return commandrun.NewRunner(time.Minute, "go", "build", "-o", query.BinaryPath, ".")
}
