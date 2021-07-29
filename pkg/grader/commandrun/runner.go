package commandrun

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/grader"
	"github.com/pkg/errors"
)

type Runner interface {
	Run(ctx context.Context, options RunOptions) (*RunResult, error)
}

type runner struct {
	timeout time.Duration
	cmdName string
	cmdArgs []string
}

var _ Runner = (*runner)(nil)

func NewRunner(timeout time.Duration, cmdName string, cmdArgs ...string) Runner {
	return &runner{
		timeout: timeout,
		cmdName: cmdName,
		cmdArgs: cmdArgs,
	}
}

func (r *runner) Run(ctx context.Context, options RunOptions) (*RunResult, error) {
	// TODO limit memory consumption
	ctx, cancel := context.WithTimeout(ctx, r.timeout*2)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, r.cmdName, r.cmdArgs...)
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	cmd.Dir = options.Dir
	cmd.Stdin = strings.NewReader(options.Stdin)

	start := time.Now()
	err := cmd.Run()
	timeToRun := time.Now().Sub(start)

	var resourceInfo *ResourceInfo
	if procState := cmd.ProcessState; procState != nil {
		resourceInfo = r.getResourceInfo(procState)
	}

	result := &RunResult{
		Stdout:       stdout.String(),
		Stderr:       stderr.String(),
		Duration:     timeToRun,
		ResourceInfo: resourceInfo,
	}

	if err == nil {
		return result, nil
	}

	exitCode, extracted := r.extractExitCode(err)
	if !extracted {
		return nil, err
	}

	result.ExitCode = exitCode

	return result, nil
}

func (r *runner) getResourceInfo(procState *os.ProcessState) *ResourceInfo {
	resourceUsage := procState.SysUsage().(*syscall.Rusage)
	return &ResourceInfo{
		Memory: r.getMemoryUsage(resourceUsage),
		Time:   r.getExecutionTime(resourceUsage),
	}

}

func (*runner) getMemoryUsage(resourceUsage *syscall.Rusage) grader.MemoryAmount {
	if runtime.GOOS == "darwin" {
		return grader.MemoryAmount(resourceUsage.Maxrss) * grader.MemoryAmountByte
	}

	return grader.MemoryAmount(resourceUsage.Maxrss) * grader.MemoryAmountKByte
}

func (*runner) getExecutionTime(resourceUsage *syscall.Rusage) time.Duration {
	return time.Duration(resourceUsage.Utime.Nano()) + time.Duration(resourceUsage.Stime.Nano())
}

func (*runner) extractExitCode(err error) (int, bool) {
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode(), true
	}

	return 0, false
}
