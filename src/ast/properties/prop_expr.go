package properties

import (
	"Falcon/ast/blockly"
	"Falcon/lex"
	"Falcon/sugar"
)

type Prop struct {
	Where lex.Token
	On    blockly.Expr
	Name  string
}

type Signature struct {
	Module    string
	BlockType string
	ValueName string
	Extras    []string
}

func makeSignature(module string, blockType string, valueName string, extras ...string) *Signature {
	return &Signature{Module: module, BlockType: blockType, ValueName: valueName, Extras: extras}
}

var properties = map[string]*Signature{
	"textLen": makeSignature("text", "text_length", "VALUE"),

	"listLen":     makeSignature("list", "lists_length", "LIST"),
	"random":      makeSignature("list", "lists_pick_random_item", "LIST"),
	"reverseList": makeSignature("list", "lists_reverse", "LIST"),
	"toCsvRow":    makeSignature("list", "lists_to_csv_row", "LIST"),
	"toCsvTable":  makeSignature("list", "lists_to_csv_table", "LIST"),
	"sort":        makeSignature("list", "lists_sort", "LIST"),
	"allButFirst": makeSignature("list", "lists_but_first", "LIST"),
	"allButLast":  makeSignature("list", "lists_but_last", "LIST"),
	"pairsToDict": makeSignature("list", "dictionaries_alist_to_dict", "PAIRS"),

	"keys":    makeSignature("dict", "dictionaries_getters", "DICT", "KEYS"),
	"values":  makeSignature("dict", "dictionaries_getters", "DICT", "VALUES"),
	"dictLen": makeSignature("dict", "dictionaries_length", "DICT"),
	"toPairs": makeSignature("dict", "dictionaries_dict_to_alist", "DICT"),
}

func (p *Prop) String() string {
	return sugar.Format("%.%", p.On.String(), p.Name)
}

func (p *Prop) Blockly() blockly.Block {
	signature, ok := properties[p.Name]
	if !ok {
		p.Where.Error("Unknown property access .%", p.Name)
	}
	switch signature.Module {
	case "text":
		return p.textProp(signature)
	case "list":
		return p.listProp(signature)
	case "dict":
		return p.dictProps(signature)
	default:
		panic("Unknown undefined module " + signature.Module)
	}
}

func (p *Prop) simpleOperand(blockType string, valueName string) blockly.Block {
	return blockly.Block{Type: blockType, Values: []blockly.Value{{Name: valueName, Block: p.On.Blockly()}}}
}
