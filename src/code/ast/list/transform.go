package list

import (
	"Falcon/code/ast/blockly"
	"Falcon/code/lex"
	"Falcon/code/sugar"
	"strconv"
	"strings"
)

type Transformer struct {
	Where       *lex.Token
	List        blockly.Expr
	Name        string
	Args        []blockly.Expr
	Names       []string
	Transformer blockly.Expr
}

type Signature struct {
	ArgSize  int
	NameSize int
}

func makeSignature(argSize int, nameSize int) *Signature {
	return &Signature{ArgSize: argSize, NameSize: nameSize}
}

var transformers = map[string]*Signature{
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
			blockly.JoinExprs(", ", t.Args),
			strings.Join(t.Names, ", "),
			t.Transformer.String())
	}
}

func (t *Transformer) Blockly() blockly.Block {
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

func (t *Transformer) Consumable() bool {
	return true
}

func (t *Transformer) max() blockly.Block {
	return blockly.Block{
		Type: "lists_maximum_value",
		Fields: []blockly.Field{
			{Name: "VAR1", Value: t.Names[0]},
			{Name: "VAR2", Value: t.Names[1]},
		},
		Values: []blockly.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "COMPARE", Block: t.Transformer.Blockly()},
		},
	}
}

func (t *Transformer) min() blockly.Block {
	return blockly.Block{
		Type: "lists_minimum_value",
		Fields: []blockly.Field{
			{Name: "VAR1", Value: t.Names[0]},
			{Name: "VAR2", Value: t.Names[1]},
		},
		Values: []blockly.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "COMPARE", Block: t.Transformer.Blockly()},
		},
	}
}

func (t *Transformer) listSortByKey() blockly.Block {
	return blockly.Block{
		Type:   "lists_sort_key",
		Fields: []blockly.Field{{Name: "VAR", Value: t.Names[0]}},
		Values: []blockly.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "KEY", Block: t.Transformer.Blockly()},
		},
	}
}

func (t *Transformer) listSort() blockly.Block {
	return blockly.Block{
		Type: "lists_sort_comparator",
		Fields: []blockly.Field{
			{Name: "VAR1", Value: t.Names[0]},
			{Name: "VAR2", Value: t.Names[1]},
		},
		Values: []blockly.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "COMPARE", Block: t.Transformer.Blockly()},
		},
	}
}

func (t *Transformer) listReduce() blockly.Block {
	return blockly.Block{
		Type: "lists_reduce",
		Fields: []blockly.Field{
			{Name: "VAR1", Value: t.Names[0]},
			{Name: "VAR2", Value: t.Names[1]},
		},
		Values: []blockly.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "INITANSWER", Block: t.Args[0].Blockly()},
			{Name: "COMBINE", Block: t.Transformer.Blockly()},
		},
	}
}

func (t *Transformer) listFilter() blockly.Block {
	return blockly.Block{
		Type:   "lists_filter",
		Fields: []blockly.Field{{Name: "VAR", Value: t.Names[0]}},
		Values: []blockly.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "TEST", Block: t.Transformer.Blockly()},
		},
	}
}

func (t *Transformer) listMap() blockly.Block {
	return blockly.Block{
		Type:   "lists_map",
		Fields: []blockly.Field{{Name: "VAR", Value: t.Names[0]}},
		Values: []blockly.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "TO", Block: t.Transformer.Blockly()},
		},
	}
}
