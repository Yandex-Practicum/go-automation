package snippetcomment_test

import (
	"testing"

	"github.com/Yandex-Practicum/go-automation/automation/gotools/pkg/snippet/snippetcomment"
	"github.com/stretchr/testify/require"
)

func TestValidateComment(t *testing.T) {
	testCases := []struct {
		Name                 string
		Content              string
		ExpectError          bool
		ExpectedErrorMessage string
	}{
		{
			Name:                 "OK",
			Content:              "это хороший коммент без точки на конце",
			ExpectError:          false,
			ExpectedErrorMessage: "",
		},
		{
			Name:                 "Empty comment",
			Content:              "",
			ExpectError:          true,
			ExpectedErrorMessage: "Empty comments are not allowed",
		},
		{
			Name:                 "Dot at the end",
			Content:              "точка на конце.",
			ExpectError:          true,
			ExpectedErrorMessage: "Do not use . at the end of line comments",
		},
		{
			Name:                 "Ellipsis at the end",
			Content:              "многоточие на конце...",
			ExpectError:          false,
			ExpectedErrorMessage: "",
		},
		{
			Name:                 "Upper case first letter",
			Content:              "Большая буква",
			ExpectError:          true,
			ExpectedErrorMessage: "First letter must be in lower case",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := snippetcomment.ValidateComment(snippetcomment.Comment{
				Content: tc.Content,
			})

			if !tc.ExpectError {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.EqualValues(t, tc.ExpectedErrorMessage, err.Error())
			}
		})
	}
}

func TestValidateDocComment(t *testing.T) {
	testCases := []struct {
		Name                 string
		Content              string
		EntityNames          []string
		ExpectError          bool
		ExpectedErrorMessage string
	}{
		{
			Name:                 "OK",
			Content:              "Foo дока.",
			EntityNames:          []string{"Foo"},
			ExpectError:          false,
			ExpectedErrorMessage: "",
		},
		{
			Name:                 "Empty comment",
			Content:              "",
			EntityNames:          []string{"Foo"},
			ExpectError:          true,
			ExpectedErrorMessage: "Empty comments are not allowed",
		},
		{
			Name:                 "Wrong prefix",
			Content:              "Bar дока",
			EntityNames:          []string{"Foo"},
			ExpectError:          true,
			ExpectedErrorMessage: "Doc comment must start with documented entity name (need prefix Foo)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := snippetcomment.ValidateDocComment(snippetcomment.DocComment{
				Content:       tc.Content,
				EntitiesNames: tc.EntityNames,
			})

			if !tc.ExpectError {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.EqualValues(t, tc.ExpectedErrorMessage, err.Error())
			}
		})
	}
}
