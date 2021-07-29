package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/grader/caserun"
	"github.com/pkg/errors"
)

var prettyOutput bool

func init() {
	flag.BoolVar(&prettyOutput, "pretty", false, "pretty print output")
}

func main() {
	flag.Parse()

	queryData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var query caserun.ComparisonQuery
	if err := json.Unmarshal(queryData, &query); err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("bad input:\"%s\"", queryData)))
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
