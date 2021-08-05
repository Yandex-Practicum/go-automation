package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
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

	fileErrors := make(map[string][]positionedError)
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
				fileErrors[path] = append(fileErrors[path], positionedError{
					Msg: fmt.Sprintf(
						"Failed to parse snippet at position %d \"%s\"", snippet.Position.Start,
						getDemonstrationPrefix(snippet.Text, 50),
					),
					Position: snippet.Position.Start,
				})
			}

			parsedSnippets = append(parsedSnippets, parsedSnippet)
		}

		for _, snippet := range parsedSnippets {
			comments := snippetcomment.ExtractComments(snippet)
			for _, comment := range comments.Comments {
				if err := snippetcomment.ValidateComment(comment); err != nil {
					fileErrors[path] = append(fileErrors[path], positionedError{
						Msg: fmt.Sprintf(
							"Wrongly formatted comment \"%s\": %s", comment.Content, err.Error(),
						),
						Position: int(comment.StartPosition),
					})
				}
			}

			for _, comment := range comments.DocComments {
				if err := snippetcomment.ValidateDocComment(comment); err != nil {
					fileErrors[path] = append(fileErrors[path], positionedError{
						Msg: fmt.Sprintf(
							"Wrongly formatted doc comment \"%s\": %s", comment.Content, err.Error(),
						),
						Position: int(comment.StartPosition),
					})
				}
			}
		}
	}

	if len(fileErrors) == 0 {
		return
	}

	log.Fatalf("Check comments in your lesson files\n %s", formatError(fileErrors))
}

type positionedError struct {
	Msg      string
	Position int
}

func formatError(errorMsgsPerFile map[string][]positionedError) string {
	files := make([]string, 0, len(errorMsgsPerFile))
	for f := range errorMsgsPerFile {
		files = append(files, f)
	}
	sort.Strings(files)

	var sb strings.Builder
	for _, fileName := range files {
		sb.WriteString(fileName)
		sb.WriteString(":\n")

		sort.Slice(errorMsgsPerFile[fileName], func(i, j int) bool {
			return errorMsgsPerFile[fileName][i].Position < errorMsgsPerFile[fileName][j].Position
		})

		for _, msg := range errorMsgsPerFile[fileName] {
			sb.WriteString("\t")
			sb.WriteString(msg.Msg)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func getDemonstrationPrefix(content string, prefixLength int) string {
	if len(content) <= prefixLength {
		return content
	}

	return content[:prefixLength] + "..."
}
