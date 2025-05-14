package types

import (
	"Falcon/sugar"
	"strconv"
)

//go:generate stringer -type=Type
type Type int

const (
	Operator Type = iota

	OpenCurve
	CloseCurve
	OpenSquare
	CloseSquare
	OpenCurly
	CloseCurly
	Equals
	Dot
	Comma
	Hyphen
	LesserThan
	LesserThanEquals
	GreaterThan
	GreaterThanEquals
	RightArrow

	Number
	Text
	Bool
	Alpha
)

type Token struct {
	Type    Type
	Content *string
	Line    int
}

func (t Token) String() string {
	if t.Content == nil {
		return t.Type.String()
	}
	return t.Type.String() + " " + *t.Content
}

func (t Token) Error(message string, args ...string) {
	panic("[line " + strconv.Itoa(t.Line) + "] " + sugar.Format(message, args...))
}
