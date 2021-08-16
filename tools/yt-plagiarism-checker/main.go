package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/plagiarismchecker"
	"github.com/caarlos0/env"
)

type Config struct {
	UserKey        string   `env:"USER_KEY"`
	Files          []string `env:"FILES"`
	Visible        bool     `env:"VISIBLE" envDefault:"false"`
	MinUniq        float32  `env:"MIN_UNIQ" envDefault:"0"`
	ExceptDomains  []string `env:"EXCEPT_DOMAINS" envDefault:""`
	RemoveSnippets bool     `env:"REMOVE_SNIPPETS" envDefault:"false"`
}

var filenameRegexp = regexp.MustCompile(`\.md$`)

func main() {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("::error ::%s", err)
		os.Exit(1)
	}

	checker := plagiarismchecker.New(cfg.UserKey, cfg.Visible, cfg.ExceptDomains, cfg.RemoveSnippets)

	uids, err := textToReview(checker, cfg.Files)
	if err != nil {
		os.Exit(2)
	}

	fmt.Println("::group::Plagiarism checker links")
	for fileName, uid := range uids {
		fmt.Printf("%s : https://text.ru/antiplagiat/%s\n", fileName, uid)
	}
	fmt.Println("::endgroup::")

	exitCode := 0
	for fileName, uid := range uids {
		uniq, err := checker.GetResult(uid)
		if err != nil {
			fmt.Printf("::error file=%s::%s\n", fileName, err)
			continue
		}
		if uniq < cfg.MinUniq {
			exitCode = 3
			fmt.Printf("::error file=%s::Text uniq is %.2f < %.2f\n", fileName, uniq, cfg.MinUniq)
		}
	}

	os.Exit(exitCode)
}

func textToReview(checker *plagiarismchecker.Checker, files []string) (map[string]string, error) {
	uids := make(map[string]string)
	for _, fileName := range files {
		if !filenameRegexp.MatchString(fileName) {
			continue
		}

		text, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Printf("::error file=%s::%s\n", fileName, err)
			return nil, err
		}

		uid, err := checker.AddText(string(text))
		if err != nil {
			fmt.Printf("::error file=%s::%s\n", fileName, err)
			continue
		}

		uids[fileName] = uid
	}

	return uids, nil
}
