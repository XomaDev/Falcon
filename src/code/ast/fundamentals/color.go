package fundamentals

import (
	"Falcon/code/ast/blockly"
	"Falcon/code/lex"
)

type Color struct {
	Where *lex.Token
	Name  string
}

func (c *Color) Yail() string {
	//TODO implement me
	panic("implement me")
}

type Signature struct {
	Code      string
	BlockType string
}

// TODO:
//
//	There are a lot more possible colors!
//	When parsing from XML blockly, it contains a hexadecimal,
//	So we do not need this table at all.
//	Invent a new way to use hex in the language, perhaps color(#ffffff) or something like that
var colorsCodes = map[string]Signature{
	"white":     {Code: "#ffffff", BlockType: "color_white"},
	"black":     {Code: "#000000", BlockType: "color_black"},
	"red":       {Code: "#ff0000", BlockType: "color_red"},
	"pink":      {Code: "#ffafaf", BlockType: "color_pink"},
	"orange":    {Code: "#ffc800", BlockType: "color_orange"},
	"yellow":    {Code: "#ffff00", BlockType: "color_yellow"},
	"green":     {Code: "#00ff00", BlockType: "color_green"},
	"cyan":      {Code: "#00ffff", BlockType: "color_cyan"},
	"blue":      {Code: "#0000ff", BlockType: "color_blue"},
	"magenta":   {Code: "#ff00ff", BlockType: "color_magenta"},
	"lightGray": {Code: "#cccccc", BlockType: "color_light_gray"},
	"gray":      {Code: "#888888", BlockType: "color_gray"},
	"darkGray":  {Code: "#444444", BlockType: "color_dark_gray"},
}

func (c *Color) String() string {
	return "color:" + c.Name
}

func (c *Color) Blockly() blockly.Block {
	signature, ok := colorsCodes[c.Name]
	if !ok {
		c.Where.Error("Unknown color name '%'", c.Name)
	}
	return blockly.Block{
		Type:   signature.BlockType,
		Fields: []blockly.Field{{Name: "COLOR", Value: signature.Code}},
	}
}

func (c *Color) Continuous() bool {
	return true
}

func (c *Color) Consumable() bool {
	return true
}
