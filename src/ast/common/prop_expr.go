package common

import (
	"Falcon/ast/blockly"
	"Falcon/lex"
	"Falcon/sugar"
)

var Properties = map[string][]string{
	"len":           {"text", "text_length", "VALUE"},
	"isEmpty":       {"text", "text_isEmpty", "VALUE"},
	"trim":          {"text", "text_trim", "TEXT"},
	"upper":         {"text", "text_changeCase", "TEXT", "UPCASE"},
	"lower":         {"text", "text_changeCase", "TEXT", "DOWNCASE"},
	"splitAtSpaces": {"text", "text_split_at_spaces", "TEXT"},
	"reverse":       {"text", "text_reverse", "VALUE"},
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
	case "text_length", "text_isEmpty", "text_trim", "text_split_at_spaces", "text_reverse":
		return p.simpleOperand(blockType, valName)
	case "text_changeCase":
		return p.textChangeCase(blockType, valName, tags[3])
	default:
		panic("Not implemented text property " + blockType)
	}
}

func (p *Prop) textChangeCase(blockType string, valName string, fieldOp string) blockly.Block {
	return blockly.Block{
		Type:   blockType,
		Fields: []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values: []blockly.Value{{Name: valName, Block: p.On.Blockly()}},
	}
}

func (p *Prop) simpleOperand(blockType string, valueName string) blockly.Block {
	return blockly.Block{Type: blockType, Values: []blockly.Value{{Name: valueName, Block: p.On.Blockly()}}}
}
