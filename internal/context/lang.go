package context

import (
	"strings"
)

type Language string

var (
	LangTR Language = "TR"
	LangEN Language = "EN"

	Languages = map[Language]bool{
		LangTR: true,
		LangEN: true,
	}

	LangDefault = LangTR
)

func ToLanguage(lang any) Language {
	if lang == nil {
		return LangDefault
	}

	var (
		l  Language
		ok bool
	)

	if l, ok = lang.(Language); !ok {
		if lstr, ok := lang.(string); ok {
			l = Language(strings.ToUpper(strings.TrimSpace(lstr)))
		}
	}

	if _, ok := Languages[l]; !ok {
		return LangDefault
	}

	return l
}

func Lang(ctx Ctx) Language {
	lang := ctx.Value(KeyLang)
	return ToLanguage(lang)
}
