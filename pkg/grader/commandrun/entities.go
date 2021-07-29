package commandrun

import (
	"time"

	"github.com/Yandex-Practicum/go-automation/pkg/grader"
)

type RunResult struct {
	ExitCode     int
	Stdout       string
	Stderr       string
	Duration     time.Duration
	ResourceInfo *ResourceInfo
}

type RunOptions struct {
	Dir   string
	Stdin string
}

type ResourceInfo struct {
	Memory grader.MemoryAmount
	Time   time.Duration
}
