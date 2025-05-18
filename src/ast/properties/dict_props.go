package properties

import blky "Falcon/ast/blockly"

func (p *Prop) dictProps(signature *Signature) blky.Block {
	switch signature.BlockType {
	case "dictionaries_length", "dictionaries_dict_to_alist":
		return p.simpleOperand(signature.BlockType, signature.ValueName)
	case "dictionaries_getters":
		return blky.Block{
			Type:   "dictionaries_getters",
			Fields: []blky.Field{{Name: "OP", Value: signature.Extras[0]}},
			Values: []blky.Value{{Name: "DICT", Block: p.On.Blockly()}},
		}
	default:
		panic("Not implemented text property " + signature.BlockType)
	}
}
