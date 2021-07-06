package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/grader/caserun"
)

var prettyOutput bool

func init() {
	flag.BoolVar(&prettyOutput, "pretty", false, "pretty print output")
}

func main() {
	flag.Parse()

	var queryData string
	if _, err := fmt.Scanln(&queryData); err != nil {
		panic(err)
	}

	var query caserun.ComparisonQuery
	if err := json.Unmarshal([]byte(queryData), &query); err != nil {
		panic(err)
	}

	comparisonRunner := caserun.NewComparisonRunner(caserun.NewRunner())
	report, err := comparisonRunner.CompareModules(context.Background(), query)
	if err != nil {
		panic(err)
	}

	resultData, err := marshalReport(report)
	if err != nil {
		panic(err)
	}

	fmt.Print(string(resultData))
}

func marshalReport(report *caserun.ComparisonReport) ([]byte, error) {
	if !prettyOutput {
		return json.Marshal(report)
	}

	return json.MarshalIndent(report, "", "  ")
}
