package ast

import "strings"

func PrimitiveCall(name string, block string, args []Expr, types string) string {
	yail := "(call-yail-primitive "
	yail += name
	yail += " (*list-for-runtime* "
	yail += JoinYailExprs(" ", args)
	yail += ") '("
	yail += types
	yail += ") \""
	yail += block
	yail += "\")"
	return yail
}

func JoinYailExprs(separator string, expressions []Expr) string {
	exprStrings := make([]string, len(expressions))
	for i, expr := range expressions {
		exprStrings[i] = expr.Yail()
	}
	return strings.Join(exprStrings, separator)
}
