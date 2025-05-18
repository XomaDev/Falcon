package properties

import "Falcon/ast/blockly"

func (p *Prop) textProp(signature *Signature) blockly.Block {
	switch signature.BlockType {
	case "text_length", "text_isEmpty", "text_trim", "text_split_at_spaces", "text_reverse", "csvRowToList", "csvTableToList":
		return p.simpleOperand(signature.BlockType, signature.ValueName)
	case "text_changeCase":
		return p.textChangeCase(signature.BlockType, signature.ValueName, signature.Extras[0])
	default:
		panic("Not implemented text property " + signature.BlockType)
	}
}

func (p *Prop) textChangeCase(blockType string, valName string, fieldOp string) blockly.Block {
	return blockly.Block{
		Type:   blockType,
		Fields: []blockly.Field{{Name: "OP", Value: fieldOp}},
		Values: []blockly.Value{{Name: valName, Block: p.On.Blockly()}},
	}
}
