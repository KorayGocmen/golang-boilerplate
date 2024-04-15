package str

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/collate"
	"golang.org/x/text/collate/build"
	"golang.org/x/text/language"
)

func SnakeToCamel(s string) string {
	var camel string
	for _, piece := range strings.Split(s, "_") {
		camel += cases.Title(language.Und, cases.NoLower).String(piece)
	}
	return camel
}

func Sorter(alphabet string) (*collate.Collator, error) {
	b := build.NewBuilder()

	for i, x := range alphabet {
		err := b.Add([]rune{x}, [][]int{{i}}, nil)
		if err != nil {
			return nil, err
		}
	}

	t, err := b.Build()
	if err != nil {
		return nil, err
	}

	return collate.NewFromTable(t), nil
}

func Compare(alphabet, s1, s2 string) int {
	for i := 0; i < len(s1) && i < len(s2); i++ {
		if strings.IndexByte(alphabet, s1[i]) < strings.IndexByte(alphabet, s2[i]) {
			return -1
		} else if strings.IndexByte(alphabet, s1[i]) > strings.IndexByte(alphabet, s2[i]) {
			return 1
		}
	}

	if len(s1) < len(s2) {
		return -1
	}

	if len(s1) > len(s2) {
		return 1
	}

	return 0
}
