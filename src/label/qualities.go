package label

import (
	"Falcon/sugar"
	"strconv"
)

//go:generate stringer -type=Quality
type Quality int

const (
	// Hard Qualities
	OpenCurve Quality = iota
	CloseCurve
	OpenSquare
	CloseSquare
	OpenCurly
	CloseCurly
	Assignment
	Dot
	Comma
	Question
	Not
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

	If
	Elif
	Else

	// Soft Qualities
	Operator
	Unary
	Equality
)

type Token struct {
	Quality      Quality
	AllQualities []Quality
	Content      *string
	Line         int
}

func (t Token) String() string {
	if t.Content == nil {
		return t.Quality.String()
	}
	return t.Quality.String() + " " + *t.Content
}

func (t Token) Error(message string, args ...string) {
	panic("[line " + strconv.Itoa(t.Line) + "] " + sugar.Format(message, args...))
}
