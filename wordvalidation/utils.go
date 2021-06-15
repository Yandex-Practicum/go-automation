package wordvalidation

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"

	"github.com/Yandex-Practicum/go-automation/filesearch"
	"github.com/pkg/errors"
)

func Validate(rootDir string) error {
	mdFiles, err := filesearch.GetFileWithExtensionPaths(rootDir, "md")
	if err != nil {
		return err
	}

	errMsgs := make([]string, 0, len(mdFiles))
	for _, mdFileName := range mdFiles {
		fileContent, err := ioutil.ReadFile(mdFileName)
		if err != nil {
			return err
		}

		bagOfWords := makeStringSet(
			splitTextIntoLexems(
				normalizeText(string(fileContent)),
				5,
			),
		)

		fileStopWords := searchStopWords(bagOfWords, makeStringSet(stopWords))

		if len(fileStopWords) > 0 {
			errMsgs = append(errMsgs, fmt.Sprintf("file %s has stop words %v", mdFileName, fileStopWords))
		}
	}

	if len(errMsgs) == 0 {
		return nil
	}

	return errors.New(
		fmt.Sprintf("stop words found:\n%s\nhttps://smsreader.github.io/stopwords/ for convenient UI", strings.Join(errMsgs, "\n")),
	)
}

func searchStopWords(bagOfWords map[string]struct{}, stopWords map[string]struct{}) []string {
	result := make([]string, 0, len(bagOfWords))
	for word := range bagOfWords {
		if _, ok := stopWords[word]; ok {
			result = append(result, word)
		}
	}

	sort.Strings(result)
	return result
}

func makeStringSet(strings []string) map[string]struct{} {
	result := make(map[string]struct{}, len(strings))
	for _, s := range strings {
		result[s] = struct{}{}
	}

	return result
}

var spaceRegexp = regexp.MustCompile("\\s+")

func splitTextIntoLexems(text string, maxLexemLength int) []string {
	words := spaceRegexp.Split(text, -1)

	result := make([]string, 0, len(words)*5)
	for start := range words {
		for end := start; end != start+maxLexemLength; end++ {
			if end >= len(words) {
				continue
			}

			result = append(result, strings.Join(words[start:end], ", "))
		}
	}

	return result
}

var punctuationSymbols = regexp.MustCompile(`[.,!?\\]+`)

func normalizeText(text string) string {
	lower := strings.ToLower(text)
	return punctuationSymbols.ReplaceAllString(lower, "")
}
