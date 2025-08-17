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

	TextEquals
	TextNotEquals

	TextLessThan
	TextGreaterThan

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
	Underscore
	At

	True
	False
	Text
	Number
	Name
	Color

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
	WalkAll
	Global
	Local
	This
	Func
)
