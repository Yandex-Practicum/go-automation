package main

import (
	"fmt"
	"github.com/Yandex-Practicum/go-automation/automation/gotools/plagiarismchecker"
	"github.com/caarlos0/env"
	"os"
	"regexp"
)

type Config struct {
	UserKey string   `env:"USER_KEY"`
	Files   []string `env:"FILES"`
	Visible bool     `env:"VISIBLE" envDefault:"false"`
	MinUniq float32  `env:"MIN_UNIQ" envDefault:"0"`
}

var filenameRegexp = regexp.MustCompile(`\.md$`)

func main() {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("::error ::%s", err)
		os.Exit(1)
	}

	checker := plagiarismchecker.New(cfg.UserKey, cfg.Visible)
	uids := make(map[string]string)
	for _, fileName := range cfg.Files {
		if !filenameRegexp.MatchString(fileName) {
			continue
		}

		text, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Printf("::error file=%s::%s\n", fileName, err)
			os.Exit(2)
		}

		uid, err := checker.AddText(string(text))
		if err != nil {
			fmt.Printf("::error file=%s::%s\n", fileName, err)
			continue
		}

		uids[fileName] = uid
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
		}
		if uniq < cfg.MinUniq {
			exitCode = 3
			fmt.Printf("::error file=%s::Text uniq is %.2f < %.2f\n", fileName, uniq, cfg.MinUniq)
		}
	}

	os.Exit(exitCode)
}
