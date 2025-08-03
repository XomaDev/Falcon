package list

import (
	blky "Falcon/ast/blockly"
	"Falcon/lex"
	"Falcon/sugar"
	"strconv"
	"strings"
)

type Transformer struct {
	Where       *lex.Token
	List        blky.Expr
	Name        string
	Args        []blky.Expr
	Names       []string
	Transformer blky.Expr
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
		pFormat := "%.% { % ->\n%}"
		if !t.List.Continuous() {
			pFormat = "(%).% { % ->\n%}"
		}
		return sugar.Format(pFormat,
			t.List.String(),
			t.Name,
			strings.Join(t.Names, ", "),
			blky.Pad(t.Transformer))
	} else {
		pFormat := "%.%(%) { % ->\n%}"
		if !t.List.Continuous() {
			pFormat = "(%).%(%) { % ->\n%}"
		}
		return sugar.Format(pFormat,
			t.List.String(),
			t.Name,
			blky.JoinExprs(", ", t.Args),
			strings.Join(t.Names, ", "),
			blky.Pad(t.Transformer))
	}
}

func (t *Transformer) Blockly() blky.Block {
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
	return false
}

func (t *Transformer) max() blky.Block {
	return blky.Block{
		Type: "lists_maximum_value",
		Fields: []blky.Field{
			{Name: "VAR1", Value: t.Names[0]},
			{Name: "VAR2", Value: t.Names[1]},
		},
		Values: []blky.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "COMPARE", Block: t.Transformer.Blockly()},
		},
		Consumable: true,
	}
}

func (t *Transformer) min() blky.Block {
	return blky.Block{
		Type: "lists_minimum_value",
		Fields: []blky.Field{
			{Name: "VAR1", Value: t.Names[0]},
			{Name: "VAR2", Value: t.Names[1]},
		},
		Values: []blky.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "COMPARE", Block: t.Transformer.Blockly()},
		},
		Consumable: true,
	}
}

func (t *Transformer) listSortByKey() blky.Block {
	return blky.Block{
		Type:   "lists_sort_key",
		Fields: []blky.Field{{Name: "VAR", Value: t.Names[0]}},
		Values: []blky.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "KEY", Block: t.Transformer.Blockly()},
		},
		Consumable: true,
	}
}

func (t *Transformer) listSort() blky.Block {
	return blky.Block{
		Type: "lists_sort_comparator",
		Fields: []blky.Field{
			{Name: "VAR1", Value: t.Names[0]},
			{Name: "VAR2", Value: t.Names[1]},
		},
		Values: []blky.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "COMPARE", Block: t.Transformer.Blockly()},
		},
		Consumable: true,
	}
}

func (t *Transformer) listReduce() blky.Block {
	return blky.Block{
		Type: "lists_reduce",
		Fields: []blky.Field{
			{Name: "VAR1", Value: t.Names[0]},
			{Name: "VAR2", Value: t.Names[1]},
		},
		Values: []blky.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "INITANSWER", Block: t.Args[0].Blockly()},
			{Name: "COMBINE", Block: t.Transformer.Blockly()},
		},
		Consumable: true,
	}
}

func (t *Transformer) listFilter() blky.Block {
	return blky.Block{
		Type:   "lists_filter",
		Fields: []blky.Field{{Name: "VAR", Value: t.Names[0]}},
		Values: []blky.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "TEST", Block: t.Transformer.Blockly()},
		},
		Consumable: true,
	}
}

func (t *Transformer) listMap() blky.Block {
	return blky.Block{
		Type:   "lists_map",
		Fields: []blky.Field{{Name: "VAR", Value: t.Names[0]}},
		Values: []blky.Value{
			{Name: "LIST", Block: t.List.Blockly()},
			{Name: "TO", Block: t.Transformer.Blockly()},
		},
		Consumable: true,
	}
}
