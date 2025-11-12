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

func CombineSignatures(first []Signature, second []Signature) []Signature {
	combined := make([]Signature, len(first)+len(second))
	copy(combined, first)
	copy(combined[len(first):], second)
	return combined
}
