package ast

import (
	"strings"
)

func Pad(code string) string {
	return "  " + strings.Replace(code, "\n", "\n  ", -1) + "\n"
}

func PadDirect(code string) string {
	return "  " + strings.Replace(code, "\n", "\n  ", -1)
}

func PadBody(blocks []Expr) string {
	var builder strings.Builder
	for _, block := range blocks {
		builder.WriteString(Pad(block.String()))
	}
	return builder.String()
}
