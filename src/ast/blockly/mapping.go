package blockly

import (
	"strconv"
	"strings"
)

func ToFields(m map[string]string) []Field {
	fields := make([]Field, 0, len(m))
	for k, v := range m {
		fields = append(fields, Field{k, v})
	}
	return fields
}

func ToValues(namePrefix string, operands []Expr) []Value {
	values := make([]Value, len(operands))
	for i, operand := range operands {
		values[i] = Value{Name: namePrefix + strconv.Itoa(i), Block: operand.Blockly()}
	}
	return values
}

func MakeValues(operands []Expr, names ...string) []Value {
	if len(operands) != len(names) {
		panic("len(operands) != len(names)")
	}
	values := make([]Value, len(operands))
	for i, operand := range operands {
		values[i] = Value{Name: names[i], Block: operand.Blockly()}
	}
	return values
}

func CreateStatement(name string, body []Expr) Statement {
	headBlock := body[0].Blockly()
	bodyLen := len(body)
	currI := 1

	for currI < bodyLen {
		aBlock := body[currI].Blockly()
		headBlock.Next = &Next{Block: &aBlock}
		currI++
	}
	return Statement{Name: name, Block: &headBlock}
}

func ToStatements(namePrefix string, bodies [][]Expr) []Statement {
	statements := make([]Statement, len(bodies))
	for i, aBody := range bodies {
		statements[i] = CreateStatement(namePrefix+strconv.Itoa(i), aBody)
	}
	return statements
}

func JoinExprs(separator string, expressions []Expr) string {
	exprStrings := make([]string, len(expressions))
	for i, expr := range expressions {
		exprStrings[i] = expr.String()
	}
	return strings.Join(exprStrings, separator)
}
