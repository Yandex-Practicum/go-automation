package caserun

import "github.com/Yandex-Practicum/go-automation/automation/gotools/grader"

type ComparisonReport struct {
	OriginalSolutionReport *SuiteReport `json:"original_solution_report"`
	SolutionReport         *SuiteReport `json:"solution_report"`
}

type SuiteReport struct {
	Cases []CaseReport `json:"cases"`
}

type CaseReport struct {
	ID             CaseID           `json:"id"`
	Tag            CaseTag          `json:"tag,omitempty"`
	TimeLimitMilli grader.TimeMilli `json:"time_limit_milli"`
	TimeUsedMilli  grader.TimeMilli `json:"time_used_milli"`
	//MemoryLimit    grader.MemoryAmount `json:"memory_limit"`
	//MemoryUsed     grader.MemoryAmount `json:"memory_used"`
	UserOutput string `json:"user_output"`
}
