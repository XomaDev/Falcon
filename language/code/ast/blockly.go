package ast

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type XmlRoot struct {
	XMLName xml.Name `xml:"xml"`
	XMLNS   string   `xml:"xmlns,attr"`
	Blocks  []Block  `xml:"block"`
}

type Block struct {
	XMLName    xml.Name    `xml:"block"`
	Type       string      `xml:"type,attr"`
	Mutation   *Mutation   `xml:"mutation,omitempty"`
	Fields     []Field     `xml:"field"`
	Values     []Value     `xml:"value"`
	Statements []Statement `xml:"statement"`
	Next       *Next       `xml:"next"`
}

type Field struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type Value struct {
	XMLName xml.Name `xml:"value"`
	Name    string   `xml:"name,attr"`
	Block   Block    `xml:"block"`
}

type Mutation struct {
	XMLName       xml.Name    `xml:"mutation"`
	ItemCount     int         `xml:"items,attr,omitempty"`
	ElseIfCount   int         `xml:"elseif,attr,omitempty"`
	ElseCount     int         `xml:"else,attr,omitempty"`
	LocalNames    []LocalName `xml:"localname"`
	Args          []Arg       `xml:"arg"`
	Key           string      `xml:"key,attr,omitempty"`
	SetOrGet      string      `xml:"set_or_get,attr,omitempty"`
	PropertyName  string      `xml:"property_name,attr,omitempty"`
	IsGeneric     bool        `xml:"is_generic,attr,omitempty"`
	ComponentType string      `xml:"component_type,attr,omitempty"`
	InstanceName  string      `xml:"instance_name,attr,omitempty"`
	EventName     string      `xml:"event_name,attr,omitempty"`
	MethodName    string      `xml:"method_name,attr,omitempty"`
	Mode          string      `xml:"mode,attr,omitempty"`
	Cofounder     string      `xml:"confounder,attr,omitempty"`
	Inline        bool        `xml:"inline,attr,omitempty"`
	Name          string      `xml:"name,attr,omitempty"`
}

type LocalName struct {
	XMLName xml.Name `xml:"localname"`
	Name    string   `xml:"name,attr"`
}

type Statement struct {
	XMLName xml.Name `xml:"statement"`
	Name    string   `xml:"name,attr"`
	Block   *Block   `xml:"block"`
}

type Next struct {
	XMLName xml.Name `xml:"next"`
	Block   *Block   `xml:"block"`
}

type Arg struct {
	Name string `xml:"name,attr"`
}

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
	if body[0].Consumable() {
		panic("Cannot include a consumable call in a body")
	}
	bodyLen := len(body)
	currI := 1

	for currI < bodyLen {
		if body[currI].Consumable() {
			panic("Cannot include a consumable call in a body")
		}
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
