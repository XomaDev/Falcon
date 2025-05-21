package properties

import "Falcon/ast/blockly"

func (p *Prop) textProp(signature *Signature) blockly.Block {
	switch signature.BlockType {
	case "text_length":
		return p.simpleOperand(signature.BlockType, signature.ValueName)
	default:
		panic("Not implemented text property " + signature.BlockType)
	}
}
