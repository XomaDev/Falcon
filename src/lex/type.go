package lex

//go:generate stringer -type=Type
type Type int

const (
	Plus Type = iota
	Dash
	Times
	Slash
	Power

	LogicOr
	LogicAnd
	BitwiseOr
	BitwiseAnd
	BitwiseXor

	Equals
	NotEquals

	LessThan
	LessThanEqual
	GreatThan
	GreaterThanEqual

	OpenCurve
	CloseCurve
	OpenSquare
	CloseSquare
	OpenCurly
	CloseCurly

	Assign
	Dot
	Comma
	Question
	Not
	Colon
	DoubleColon
	RightArrow

	True
	False
	Text
	Number
	Name

	If
	Elif
	Else
	For
	To
	By
	Each
	In
	While
	Do
	Break
)
