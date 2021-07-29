package caserun

import "github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/grader"

type CaseID string

type SuiteID string

type CaseTag string

type Suite struct {
	ID    SuiteID `json:"id"`
	Cases []Case  `json:"cases"`
}

type Case struct {
	ID             CaseID           `json:"id"`
	Tag            CaseTag          `json:"tag,omitempty"`
	Input          string           `json:"input"`
	TimeLimitMilli grader.TimeMilli `json:"time_limit_milli"`
	// TODO add memory limit
}
