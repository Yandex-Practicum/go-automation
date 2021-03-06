package snippetcomment_test

import (
	"go/token"
	"strconv"
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
			ExpectedComments: makeSimpleGoDoComments("This is a doc comment", 59, "File"),
		},
		{
			Name: "AliasDocComment",
			Input: `
// This is a doc comment
type A = a
`,
			ExpectedComments: makeSimpleGoDoComments("This is a doc comment", 51, "A"),
		},
		{
			Name: "FuncDocComment",
			Input: `
// This is a doc comment
func Foo() {}
`,
			ExpectedComments: makeSimpleGoDoComments("This is a doc comment", 54, "Foo"),
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
			ExpectedComments: makeSimpleGoDoComments("This is a doc comment", 50, "A"),
		},
		{
			Name: "MultilineDocComment",
			Input: `
// This is a doc comment
// multiline one
var A int
`,
			ExpectedComments: makeSimpleGoDoComments("This is a doc comment\nmultiline one", 67, "A"),
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
						StartPosition: 43,
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
						StartPosition: 60,
					},
				},
			},
		},
		{
			Name: "GenerateComment",
			Input: `
//go:generate msgp
type AccountBalance struct { }
`,
			ExpectedComments: snippetcomment.Comments{
				DocComments: []snippetcomment.DocComment{
					{
						StartPosition: 65,
						Content:       "",
						EntitiesNames: []string{"AccountBalance"},
						IsDirective:   true,
					},
				},
			},
		},
		{
			Name: "LinterComment",
			Input: `
//lint:file-ignore U1000 ???????????????????? ???? ???????????????????????? ??????, ?????? ?????? ???? ????????????????????????
`,
			ExpectedComments: snippetcomment.Comments{
				Comments: []snippetcomment.Comment{
					{
						StartPosition: 143,
						Content:       "",
						IsDirective:   true,
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

			t.Run("SimpleComments", func(t *testing.T) {
				require.Len(t, comments.Comments, len(tc.ExpectedComments.Comments))

				for i, c := range tc.ExpectedComments.Comments {
					t.Run(strconv.Itoa(i), func(t *testing.T) {
						require.EqualValues(t, c.StartPosition, comments.Comments[i].StartPosition)
						require.EqualValues(t, c.Content, comments.Comments[i].Content)
						require.EqualValues(t, c.IsDirective, comments.Comments[i].IsDirective)
					})
				}
			})

			t.Run("DocComments", func(t *testing.T) {
				require.Len(t, comments.DocComments, len(tc.ExpectedComments.DocComments))

				for i, c := range tc.ExpectedComments.DocComments {
					require.EqualValues(t, c.StartPosition, comments.DocComments[i].StartPosition)
					require.EqualValues(t, c.Content, comments.DocComments[i].Content)
					require.EqualValues(t, c.IsDirective, comments.DocComments[i].IsDirective)
				}
			})
		})
	}
}
