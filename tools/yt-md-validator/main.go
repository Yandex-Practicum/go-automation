package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/filesearch"
	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/markdown/mdvalidation"
	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/markdown/mdvalidation/mdvalidatorblanklines"
	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/markdown/mdvalidation/mdvalidatorheaderlevels"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

var rootDir string

func init() {
	flag.StringVar(&rootDir, "root", ".", "root")
}

var allValidators = []mdvalidation.Validator{
	mdvalidatorblanklines.NewValidator(),
	mdvalidatorheaderlevels.NewValidator(),
}

func main() {
	flag.Parse()

	paths, err := filesearch.GetFileWithExtensionPaths(rootDir, "md")
	if err != nil {
		log.Panicf("failed to read markdown files: %s", err)
	}

	p := goldmark.DefaultParser()

	infosPerFile := make(map[string][]*mdvalidation.ValidationInfo)
	for _, path := range paths {
		fileContent, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read file %s: %s", fileContent, err)
		}

		fileAST := p.Parse(text.NewReader(fileContent))

		for _, v := range allValidators {
			infos, err := v.ValidateTree(context.Background(), fileAST, fileContent)
			if err != nil {
				log.Fatalf(err.Error())
			}

			if len(infos) > 0 {
				infosPerFile[path] = append(infosPerFile[path], infos...)
			}
		}
	}

	if len(infosPerFile) > 0 {
		log.Fatalf("Markdown lessons are poorly formatted: %v", formatError(infosPerFile))
	}
}

func formatError(infosPerFile map[string][]*mdvalidation.ValidationInfo) string {
	files := make([]string, 0, len(infosPerFile))
	for f := range infosPerFile {
		files = append(files, f)
	}
	sort.Strings(files)

	var sb strings.Builder
	for _, fileName := range files {
		sb.WriteString(fileName)
		sb.WriteString(":\n")

		sort.Slice(infosPerFile[fileName], func(i, j int) bool {
			return mdvalidation.GetNodeStart(infosPerFile[fileName][i].Node) < mdvalidation.GetNodeStart(infosPerFile[fileName][j].Node)
		})

		for _, info := range infosPerFile[fileName] {
			sb.WriteString("\t")
			sb.WriteString(info.Message)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
