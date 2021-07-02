package commandrun

import (
	"time"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/grader"
)

type RunResult struct {
	Stdout       string
	Stderr       string
	Duration     time.Duration
	ResourceInfo *ResourceInfo
}

type ResourceInfo struct {
	Memory grader.MemoryAmount
	Time   time.Duration
}
