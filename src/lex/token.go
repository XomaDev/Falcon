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
		return "(" + t.Type.String() + ")"
	}
	return *t.Content
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
	panic("[line " + strconv.Itoa(t.Line) + "] " + sugar.Format(message, args...))
}

type StaticToken struct {
	Type  Type
	Flags []Flag
}

func staticOf(t Type, flags ...Flag) StaticToken {
	return StaticToken{t, flags}
}

func (s *StaticToken) normal(line int, content ...string) Token {
	if len(content) > 1 {
		panic("Too many contents...")
	}
	return Token{Line: line, Type: s.Type, Flags: s.Flags, Content: &content[0]}
}
