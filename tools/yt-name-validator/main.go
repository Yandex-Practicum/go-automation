package main

import (
	"flag"

	"github.com/Yandex-Practicum/go-automation/namevalidation"
)

var rootDir string

func init() {
	flag.StringVar(&rootDir, "root", ".", "defines directory in which check file names")
}

func main() {
	flag.Parse()

	if err := namevalidation.Validate(rootDir); err != nil {
		panic(err)
	}
}
