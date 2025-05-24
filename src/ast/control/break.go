package control

import "Falcon/ast/blockly"

type Break struct {
	// Hola Amigo!
}

func (b *Break) String() string {
	return "break"
}

func (b *Break) Blockly() blockly.Block {
	return blockly.Block{Type: "controls_break", Consumable: false}
}
