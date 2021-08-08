package snippetcomment_test

import (
	"go/token"
	"testing"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetcomment"
	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetparse"
	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetsearch"
	"github.com/stretchr/testify/require"
)

func TestCommentsExtraction(t *testing.T) {
	makeSimpleGoDoComments := func(content string, pos int, entitiesNames ...string) snippetcomment.Comments {
		return snippetcomment.Comments{
			DocComments: []snippetcomment.DocComment{
				{
					Content:       content,
					StartPosition: token.Pos(pos),
					EntitiesNames: entitiesNames,
				},
			},
		}
	}

	testCases := []struct {
		Name             string
		Input            string
		ExpectedComments snippetcomment.Comments
	}{
		{
			Name:             "NoComments",
			Input:            "b.buf = b.buf[:]",
			ExpectedComments: snippetcomment.Comments{},
		},
		{
			Name: "TypeDocComment",
			Input: `
// This is a doc comment
type File struct{}
`,
			ExpectedComments: makeSimpleGoDoComments("This is a doc comment", 58, "File"),
		},
		{
			Name: "AliasDocComment",
			Input: `
// This is a doc comment
type A = a
`,
			ExpectedComments: makeSimpleGoDoComments("This is a doc comment", 50, "A"),
		},
		{
			Name: "FuncDocComment",
			Input: `
// This is a doc comment
func Foo() {}
`,
			ExpectedComments: makeSimpleGoDoComments("This is a doc comment", 53, "Foo"),
		},
		{
			Name: "PackageDocComment",
			Input: `
// This is a doc comment
package main
				`,
			ExpectedComments: makeSimpleGoDoComments("This is a doc comment", 2, "main"),
		},
		{
			Name: "VariableDocComment",
			Input: `
// This is a doc comment
var A int
`,
			ExpectedComments: makeSimpleGoDoComments("This is a doc comment", 49, "A"),
		},
		{
			Name: "MultilineDocComment",
			Input: `
// This is a doc comment
// multiline one
var A int
`,
			ExpectedComments: makeSimpleGoDoComments("This is a doc comment\nmultiline one", 66, "A"),
		},
		{
			Name: "SimpleComment",
			Input: `
// This is a simple comment
`,
			ExpectedComments: snippetcomment.Comments{
				Comments: []snippetcomment.Comment{
					{
						Content:       "This is a simple comment",
						StartPosition: 42,
					},
				},
			},
		},
		{
			Name: "SimpleMultilineComment",
			Input: `
// This is a simple comment
// multiline one
`,
			ExpectedComments: snippetcomment.Comments{
				Comments: []snippetcomment.Comment{
					{
						Content:       "This is a simple comment\nmultiline one",
						StartPosition: 59,
					},
				},
			},
		},
		{
			Name: "RealWorldExample",
			Input: `
// package
package main

// type
type Type struct{}

// func
func Foo(t Type) {}

func main() {
	// simple
	Foo(Type{})
}
`,
			ExpectedComments: snippetcomment.Comments{
				DocComments: []snippetcomment.DocComment{
					{
						Content:       "package",
						EntitiesNames: []string{"main"},
						StartPosition: 2,
					},
					{
						Content:       "type",
						EntitiesNames: []string{"Type"},
						StartPosition: 27,
					},
					{
						Content:       "func",
						EntitiesNames: []string{"Foo"},
						StartPosition: 55,
					},
				},
				Comments: []snippetcomment.Comment{
					{
						Content:       "simple",
						StartPosition: 99,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			parsedSnippet, err := snippetparse.ParseSnippet(snippetsearch.Snippet{
				Text: tc.Input,
			})
			require.NoError(t, err)

			comments := snippetcomment.ExtractComments(parsedSnippet)

			require.EqualValues(t, tc.ExpectedComments, comments)
		})
	}
}
