package blockly

import (
	"strconv"
	"strings"
)

func FieldsFromMap(m map[string]string) []Field {
	fields := make([]Field, 0, len(m))
	for k, v := range m {
		fields = append(fields, Field{k, v})
	}
	return fields
}

func ToFields(prefix string, values []string) []Field {
	fields := make([]Field, len(values))
	for i, value := range values {
		fields[i] = Field{prefix + strconv.Itoa(i), value}
	}
	return fields
}

func ToArgs(names []string) []Arg {
	args := make([]Arg, len(names))
	for i, name := range names {
		args[i] = Arg{Name: name}
	}
	return args
}

func ValuesByPrefix(namePrefix string, operands []Expr) []Value {
	values := make([]Value, len(operands))
	for i, operand := range operands {
		values[i] = Value{Name: namePrefix + strconv.Itoa(i), Block: operand.Blockly()}
	}
	return values
}

func ValueArgsByPrefix(on Expr, onName string, namePrefix string, operands []Expr) []Value {
	values := make([]Value, len(operands)+1)
	values[0] = Value{Name: onName, Block: on.Blockly()}
	for i, operand := range operands {
		values[i+1] = Value{Name: namePrefix + strconv.Itoa(i), Block: operand.Blockly()}
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

func MakeValueArgs(on Expr, onName string, operands []Expr, names ...string) []Value {
	if len(operands) != len(names) {
		panic("len(operands) != len(names)")
	}
	values := make([]Value, len(operands)+1)
	values[0] = Value{Name: onName, Block: on.Blockly()}
	for i, operand := range operands {
		values[i+1] = Value{Name: names[i], Block: operand.Blockly()}
	}
	return values
}

func CreateStatement(name string, body []Expr) Statement {
	headBlock := body[0].Blockly()
	//if headBlock.Consumable {
	//	panic("Cannot include a consumable call in a body")
	//}
	bodyLen := len(body)
	currI := 1

	for currI < bodyLen {
		aBlock := body[currI].Blockly()
		if aBlock.Consumable {
			panic("Cannot include a consumable call in a body")
		}
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

func MakeLocalNames(names ...string) []LocalName {
	localNames := make([]LocalName, len(names))
	for i, name := range names {
		localNames[i] = LocalName{Name: name}
	}
	return localNames
}

func JoinExprs(separator string, expressions []Expr) string {
	exprStrings := make([]string, len(expressions))
	for i, expr := range expressions {
		exprStrings[i] = expr.String()
	}
	return strings.Join(exprStrings, separator)
}
