package control

import (
	"Falcon/code/ast/blockly"
)

type Break struct {
	// Hola Amigo!
}

func (b *Break) String() string {
	return "break"
}

func (b *Break) Blockly() blockly.Block {
	return blockly.Block{Type: "controls_break"}
}

func (b *Break) Continuous() bool {
	return true
}

func (b *Break) Consumable() bool {
	return false
}
