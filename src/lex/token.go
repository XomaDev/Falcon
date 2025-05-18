package lex

import (
	"Falcon/sugar"
	"strconv"
)

type Token struct {
	Line    int
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

func (t *Token) HasFlag(flag Flag) bool {
	for _, f := range t.Flags {
		if f == flag {
			return true
		}
	}
	return false
}

func (t *Token) Error(message string, args ...string) {
	panic("[line " + strconv.Itoa(t.Line) + "] " + t.String() + " " + sugar.Format(message, args...))
}

type StaticToken struct {
	Type  Type
	Flags []Flag
}

func staticOf(t Type, flags ...Flag) StaticToken {
	return StaticToken{t, flags}
}

func (s *StaticToken) normal(line int, optionalContent ...string) Token {
	if len(optionalContent) > 1 {
		panic("Too many contents...")
	}
	var content string
	if len(optionalContent) == 1 {
		content = optionalContent[0]
	}
	return Token{Line: line, Type: s.Type, Flags: s.Flags, Content: &content}
}
