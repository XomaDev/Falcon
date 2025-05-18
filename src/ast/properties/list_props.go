package properties

import "Falcon/ast/blockly"

func (p *Prop) listProp(signature *Signature) blockly.Block {
	switch signature.BlockType {
	case "lists_length",
		"lists_pick_random_item",
		"lists_reverse",
		"lists_to_csv_row",
		"lists_to_csv_table",
		"lists_sort",
		"lists_but_first",
		"lists_but_last":
		return p.simpleOperand(signature.BlockType, signature.ValueName)
	default:
		panic("Not implemented list property " + signature.BlockType)
	}
}
