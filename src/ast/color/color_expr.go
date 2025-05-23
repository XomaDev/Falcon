package color

import (
	"Falcon/ast/blockly"
	"Falcon/lex"
)

type Color struct {
	Where lex.Token
	Name  string
}

type Signature struct {
	Code      string
	BlockType string
}

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
