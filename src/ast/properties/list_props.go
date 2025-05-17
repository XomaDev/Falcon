package properties

import "Falcon/ast/blockly"

func (p *Prop) listProp(signature *Signature) blockly.Block {
	switch signature.BlockType {
	case "lists_length", "lists_pick_random_item":
		return p.simpleOperand(signature.BlockType, signature.ValueName)
	default:
		panic("Not implemented list property " + signature.BlockType)
	}
}
