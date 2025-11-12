package ast

type Signature int

const (
	SignBool Signature = iota
	SignNumb
	SignText
	SignList
	SignDict
	SignComponent
	SignHelper
	SignAny
	SignVoid
)
