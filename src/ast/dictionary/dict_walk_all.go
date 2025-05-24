package dictionary

import "Falcon/ast/blockly"

type WalkAll struct {
}

func (w *WalkAll) String() string {
	return "walkAll"
}

func (w *WalkAll) Blockly() blockly.Block {
	return blockly.Block{Type: "dictionaries_walk_all", Consumable: true}
}
