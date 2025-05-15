package lex

//go:generate stringer -type=Flag
type Flag int

const (
	Operator Flag = iota
	LLogicOr
	LLogicAnd
	BBitwiseOr
	BBitwiseAnd
	BBitwiseXor

	Relational
	Equality
	Binary
	BinaryL1
	BinaryL2
	Unary

	Value
	ConstantValue

	PreserveOrder
)
