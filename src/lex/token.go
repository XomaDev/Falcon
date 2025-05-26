package lex

import (
	"Falcon/context"
	"Falcon/sugar"
	"strconv"
)

type Token struct {
	Column  int
	Row     int
	Context *context.CodeContext

	Type    Type
	Flags   []Flag
	Content *string
}

func (t *Token) String() string {
	if t.Content == nil {
		return sugar.Format("(%)", t.Type.String())
	}
	return sugar.Format("(% %)", t.Type.String(), *t.Content)
}

func (t *Token) Debug() string {
	if t.Content == nil {
		return sugar.Format("(%:% %)", strconv.Itoa(t.Row), strconv.Itoa(t.Column), t.Type.String())
	}
	return sugar.Format("(%:% % %)", strconv.Itoa(t.Row), strconv.Itoa(t.Column), t.Type.String(), *t.Content)
}

func (t *Token) HasFlag(flag Flag) bool {
	for _, f := range t.Flags {
		if f == flag {
			return true
		}
	}
	return false
}

func (t *Token) Error(message string, args ...string) {
	var wordSize = 6
	if t.Content != nil {
		wordSize = max(wordSize, len(*t.Content))
	}
	t.Context.ReportError(t.Column, t.Row, wordSize, message, args...)
}

type StaticToken struct {
	Type  Type
	Flags []Flag
}

func staticOf(t Type, flags ...Flag) StaticToken {
	return StaticToken{t, flags}
}

func (s *StaticToken) normal(
	column int,
	row int,
	ctx *context.CodeContext,
	optionalContent ...string,
) *Token {
	if len(optionalContent) > 1 {
		panic("Too many contents...")
	}
	var content string
	if len(optionalContent) == 1 {
		content = optionalContent[0]
	}
	return &Token{
		Column:  column,
		Row:     row,
		Context: ctx,

		Type:    s.Type,
		Flags:   s.Flags,
		Content: &content,
	}
}
