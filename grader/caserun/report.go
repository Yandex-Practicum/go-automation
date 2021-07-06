package caserun

import "github.com/Yandex-Practicum/go-automation/automation/gotools/grader"

type ComparisonReport struct {
	OriginalSolutionReport *SuiteReport
	SolutionReport         *SuiteReport
}

type SuiteReport struct {
	Cases []CaseReport
}

type CaseReport struct {
	ID             CaseID           `json:"id"`
	Tag            CaseTag          `json:"tag,omitempty"`
	TimeLimitMilli grader.TimeMilli `json:"timeLimitMilli"`
	TimeUsed       grader.TimeMilli `json:"timeUsed"`
	//MemoryLimit    grader.MemoryAmount `json:"memoryLimit"`
	//MemoryUsed     grader.MemoryAmount `json:"memoryUsed"`
	UserOutput string `json:"userOutput"`
}
