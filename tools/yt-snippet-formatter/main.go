package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/Yandex-Practicum/go-automation/filesearch"
	"github.com/Yandex-Practicum/go-automation/snippetformat"
)

const (
	modeValidation = "validation"
	modeFormat     = "format"
)

var (
	rootDir string
	mode    string
)

func init() {
	flag.StringVar(&rootDir, "root", ".", "defines directory in which markdown with snippets must be formatted")
	flag.StringVar(&mode, "mode", "format", "\"validation\" to check file format; \"format\" to format file")
}

func main() {
	flag.Parse()

	switch mode {
	case modeValidation, modeFormat:
	default:
		log.Panicf("unexpected mode %s", mode)
	}

	paths, err := filesearch.GetFileWithExtensionPaths(rootDir, "md")
	if err != nil {
		log.Panicf("failed to read markdown files: %s", err)
	}

	var notFormattedFiles []string
	for _, path := range paths {
		fileContent, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read file %s: %s", fileContent, err)
		}

		formattedContent, err := snippetformat.FormatText(string(fileContent))
		if err != nil {
			log.Fatalf("failed to format file %s: %s", path, err)
		}

		switch mode {
		case modeFormat:
			if err := ioutil.WriteFile(path, []byte(formattedContent), os.ModeExclusive); err != nil {
				log.Fatalf("failed to write formatted file %s: %s", path, err)
			}

		case modeValidation:
			if string(fileContent) != formattedContent {
				notFormattedFiles = append(notFormattedFiles, path)
			}
		}
	}

	if len(notFormattedFiles) > 0 {
		log.Fatalf("snippets in files %s are not properly formatted; use yt-snippet-formatter on them", notFormattedFiles)
	}
}
