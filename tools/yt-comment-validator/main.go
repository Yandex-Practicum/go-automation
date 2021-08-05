package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/filesearch"
	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetcomment"
	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetparse"
	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetsearch"
)

var (
	rootDir string
)

func init() {
	flag.StringVar(&rootDir, "root", ".", "defines directory in which markdown with snippets must be formatted")
}

func main() {
	flag.Parse()

	paths, err := filesearch.GetFileWithExtensionPaths(rootDir, "md")
	if err != nil {
		log.Panicf("failed to read markdown files: %s", err)
	}

	errorMsgsPerFile := make(map[string][]string)
	for _, path := range paths {
		fileContent, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read file %s: %s", fileContent, err)
		}

		snippets := snippetsearch.FindSnippets(string(fileContent))

		parsedSnippets := make([]snippetparse.ParsedSnippet, 0, len(snippets))
		for _, snippet := range snippets {
			parsedSnippet, err := snippetparse.ParseSnippet(snippet)
			if err != nil {
				errorMsgsPerFile[path] = append(errorMsgsPerFile[path], fmt.Sprintf(
					"Failed to parse snippet at position %d \"%s\"", snippet.Position.Start,
					getDemonstrationPrefix(snippet.Text, 50),
				))
			}

			parsedSnippets = append(parsedSnippets, parsedSnippet)
		}

		for _, snippet := range parsedSnippets {
			comments := snippetcomment.ExtractComments(snippet)
			for _, comment := range comments.Comments {
				if err := snippetcomment.ValidateComment(comment); err != nil {
					errorMsgsPerFile[path] = append(errorMsgsPerFile[path], fmt.Sprintf(
						"Wrongly formatted comment \"%s\": %s", comment.Content, err.Error(),
					))
				}
			}

			for _, comment := range comments.DocComments {
				if err := snippetcomment.ValidateDocComment(comment); err != nil {
					errorMsgsPerFile[path] = append(errorMsgsPerFile[path], fmt.Sprintf(
						"Wrongly formatted doc comment \"%s\": %s", comment.Content, err.Error(),
					))
				}
			}
		}
	}

	if len(errorMsgsPerFile) == 0 {
		return
	}

	var sb strings.Builder
	for fileName, errMsgs := range errorMsgsPerFile {
		sb.WriteString(fileName)
		sb.WriteString(":\n")
		for _, msg := range errMsgs {
			sb.WriteString("\t")
			sb.WriteString(msg)
			sb.WriteString("\n")
		}
	}

	log.Fatalf("Check comments in your lesson files\n %s", sb.String()) // TODO make beautiful formatting
}

func getDemonstrationPrefix(content string, prefixLength int) string {
	if len(content) <= prefixLength {
		return content
	}

	return content[:prefixLength] + "..."
}
