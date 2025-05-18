package list

import (
	blky "Falcon/ast/blockly"
	"Falcon/lex"
	"Falcon/sugar"
	"strconv"
	"strings"
)

type Transformer struct {
	Where       lex.Token
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
	"map":    makeSignature(0, 1),
	"filter": makeSignature(0, 1),
}

func (t *Transformer) String() string {
	if len(t.Args) == 0 {
		return sugar.Format("%.% { % ->\n%}",
			t.List.String(),
			t.Name,
			strings.Join(t.Names, ", "),
			blky.Pad(t.Transformer))
	} else {
		return sugar.Format("%.% { % ->\n%}",
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
	default:
		panic("Unimplemented list transformer! " + t.Name)
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
	}
}
