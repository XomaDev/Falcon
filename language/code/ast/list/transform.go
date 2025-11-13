package list

import (
	"Falcon/code/ast"
	"Falcon/code/lex"
	"Falcon/code/sugar"
	"strconv"
	"strings"
)

type Transformer struct {
	Where       *lex.Token
	List        ast.Expr
	Name        string
	Args        []ast.Expr
	Names       []string
	Transformer ast.Expr
}

func (t *Transformer) Yail() string {
	//TODO implement me
	panic("implement me")
}

type transformerSignature struct {
	ArgSize  int
	NameSize int
}

func makeSignature(argSize int, nameSize int) *transformerSignature {
	return &transformerSignature{ArgSize: argSize, NameSize: nameSize}
}

var transformers = map[string]*transformerSignature{
	"map":       makeSignature(0, 1),
	"filter":    makeSignature(0, 1),
	"reduce":    makeSignature(1, 2),
	"sort":      makeSignature(0, 2),
	"sortByKey": makeSignature(0, 1),
	"min":       makeSignature(0, 2),
	"max":       makeSignature(0, 2),
}

func (t *Transformer) String() string {
	if len(t.Args) == 0 {
		pFormat := "%.% { % -> % }"
		if !t.List.Continuous() {
			pFormat = "(%).% { % -> %} "
		}
		return sugar.Format(pFormat,
			t.List.String(),
			t.Name,
			strings.Join(t.Names, ", "),
			t.Transformer.String())
	} else {
		pFormat := "%.%(%) { % -> % }"
		if !t.List.Continuous() {
			pFormat = "(%).%(%) { % -> % }"
		}
		return sugar.Format(pFormat,
			t.List.String(),
			t.Name,
			ast.JoinExprs(", ", t.Args),
			strings.Join(t.Names, ", "),
			t.Transformer.String())
	}
}

func (t *Transformer) Blockly(flags ...bool) ast.Block {
	signature, ok := transformers[t.Name]
	if !ok {
		t.Where.Error("Unknown transformer '%'", t.Name)
		panic("Unreachable")
	}
	gotArgs := len(t.Args)
	if signature.ArgSize != gotArgs {
		t.Where.Error("Expected % args but got % for transformer ::%",
			strconv.Itoa(signature.ArgSize), strconv.Itoa(gotArgs), t.Name)
	}
	gotNamesLen := len(t.Names)
	if signature.NameSize != gotNamesLen {
		t.Where.Error("Expected % names but got % for transformer ::%",
			strconv.Itoa(signature.NameSize), strconv.Itoa(gotNamesLen), t.Name)
	}
	switch t.Name {
	case "map":
		return t.listMap()
	case "filter":
		return t.listFilter()
	case "reduce":
		return t.listReduce()
	case "sort":
		return t.listSort()
	case "sortByKey":
		return t.listSortByKey()
	case "min":
		return t.min()
	case "max":
		return t.max()
	default:
		panic("Unimplemented list transformer! " + t.Name)
	}
}

func (t *Transformer) Continuous() bool {
	return true
}

func (t *Transformer) Consumable(flags ...bool) bool {
	return true
}

func (t *Transformer) Signature() []ast.Signature {
	if t.Name == "min" || t.Name == "max" || t.Name == "reduce" {
		return t.Transformer.Signature()
	}
	return []ast.Signature{ast.SignList}
}

func (t *Transformer) max() ast.Block {
	return ast.Block{
		Type: "lists_maximum_value",
		Fields: []ast.Field{
			{Name: "VAR1", Value: t.Names[0]},
			{Name: "VAR2", Value: t.Names[1]},
		},
		Values: []ast.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "COMPARE", Block: t.Transformer.Blockly()},
		},
	}
}

func (t *Transformer) min() ast.Block {
	return ast.Block{
		Type: "lists_minimum_value",
		Fields: []ast.Field{
			{Name: "VAR1", Value: t.Names[0]},
			{Name: "VAR2", Value: t.Names[1]},
		},
		Values: []ast.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "COMPARE", Block: t.Transformer.Blockly()},
		},
	}
}

func (t *Transformer) listSortByKey() ast.Block {
	return ast.Block{
		Type:   "lists_sort_key",
		Fields: []ast.Field{{Name: "VAR", Value: t.Names[0]}},
		Values: []ast.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "KEY", Block: t.Transformer.Blockly()},
		},
	}
}

func (t *Transformer) listSort() ast.Block {
	return ast.Block{
		Type: "lists_sort_comparator",
		Fields: []ast.Field{
			{Name: "VAR1", Value: t.Names[0]},
			{Name: "VAR2", Value: t.Names[1]},
		},
		Values: []ast.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "COMPARE", Block: t.Transformer.Blockly()},
		},
	}
}

func (t *Transformer) listReduce() ast.Block {
	return ast.Block{
		Type: "lists_reduce",
		Fields: []ast.Field{
			{Name: "VAR1", Value: t.Names[0]},
			{Name: "VAR2", Value: t.Names[1]},
		},
		Values: []ast.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "INITANSWER", Block: t.Args[0].Blockly()},
			{Name: "COMBINE", Block: t.Transformer.Blockly()},
		},
	}
}

func (t *Transformer) listFilter() ast.Block {
	return ast.Block{
		Type:   "lists_filter",
		Fields: []ast.Field{{Name: "VAR", Value: t.Names[0]}},
		Values: []ast.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "TEST", Block: t.Transformer.Blockly()},
		},
	}
}

func (t *Transformer) listMap() ast.Block {
	return ast.Block{
		Type:   "lists_map",
		Fields: []ast.Field{{Name: "VAR", Value: t.Names[0]}},
		Values: []ast.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "TO", Block: t.Transformer.Blockly()},
		},
	}
}
