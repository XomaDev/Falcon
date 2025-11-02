package lex

import (
	"Falcon/lang/code/context"
	"Falcon/lang/code/sugar"
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
	return sugar.Format("(% %)", t.Type.String(), *t.Content)
}

func (t *Token) Debug() string {
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
	if t.Context != nil {
		(*t.Context).ReportError(t.Column, t.Row, len(*t.Content), message, args...)
	} else {
		panic(sugar.Format(message, args...))
	}
}

type StaticToken struct {
	Type  Type
	Flags []Flag
}

func staticOf(t Type, flags ...Flag) StaticToken {
	return StaticToken{t, flags}
}

func (s *StaticToken) Normal(
	column int,
	row int,
	ctx *context.CodeContext,
	content string,
) *Token {
	return &Token{
		Column:  column,
		Row:     row,
		Context: ctx,

		Type:    s.Type,
		Flags:   s.Flags,
		Content: &content,
	}
}
