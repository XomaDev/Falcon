package variables

import blky "Falcon/ast/blockly"

type Global struct {
	Name  string
	Value blky.Expr
}

func (g *Global) String() string {
	return "glob " + g.Name + " = " + g.Value.String()
}

func (g *Global) Blockly() blky.Block {
	return blky.Block{
		Type:       "global_declaration",
		Fields:     []blky.Field{{Name: "NAME", Value: g.Name}},
		Values:     []blky.Value{{Name: "VALUE", Block: g.Value.Blockly()}},
		Consumable: false,
	}
}
