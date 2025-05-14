package lex

//go:generate stringer -type=Type
type Type int

const (
	// +-*/
	Plus Type = iota
	Dash
	Times
	Slash
	Power

	// || && | & ~
	LogicOr
	LogicAnd
	BitwiseOr
	BitwiseAnd
	BitwiseXor

	// < <= > >=
	LessThan
	LessThanEqual
	GreatThan
	GreaterThanEqual

	// ()[]{}
	OpenCurve
	CloseCurve
	OpenSquare
	CloseSquare
	OpenCurly
	CloseCurly

	// =.,?!
	Equal
	Dot
	Comma
	Question
	Not
	RightArrow

	True
	False
	Text
	Number
	Name

	If
	Elif
	Else
)
