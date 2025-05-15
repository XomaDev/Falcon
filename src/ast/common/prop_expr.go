package common

import (
	"Falcon/ast/blockly"
	"Falcon/lex"
	"Falcon/sugar"
)

var Properties = map[string][]string{
	"len":     {"text", "text_length", "VALUE"},
	"isEmpty": {"text", "text_isEmpty", "VALUE"},
	"trim":    {"text", "text_trim", "TEXT"},
}

type Prop struct {
	Where lex.Token
	On    blockly.Expr
	Name  string
}

func (p *Prop) String() string {
	return sugar.Format("%.%", p.On.String(), p.Name)
}

func (p *Prop) Blockly() blockly.Block {
	tags, ok := Properties[p.Name]
	if !ok {
		p.Where.Error("Unknown property access .%", p.Name)
	}
	switch tags[0] {
	case "text":
		return p.textProp(tags)
	default:
		panic("Unknown undefined module " + tags[0])
	}
}

func (p *Prop) textProp(tags []string) blockly.Block {
	blockType := tags[1]
	valName := tags[2]
	switch blockType {
	case "text_length", "text_isEmpty", "text_trim":
		return p.simpleOperand(blockType, valName)
	default:
		panic("Not implemented text property " + blockType)
	}
}

func (p *Prop) simpleOperand(blockType string, valueName string) blockly.Block {
	return blockly.Block{Type: blockType, Values: []blockly.Value{{Name: valueName, Block: p.On.Blockly()}}}
}
