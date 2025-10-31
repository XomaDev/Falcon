package control

import (
	"Falcon/code/ast"
)

type Break struct {
	// Hola Amigo!
}

func (b *Break) Yail() string {
	//TODO implement me
	panic("implement me")
}

func (b *Break) String() string {
	return "break"
}

func (b *Break) Blockly() ast.Block {
	return ast.Block{Type: "controls_break"}
}

func (b *Break) Continuous() bool {
	return true
}

func (b *Break) Consumable() bool {
	return false
}
