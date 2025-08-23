package blockly

import (
	"encoding/xml"
)

type XmlRoot struct {
	XMLName xml.Name `xml:"xml"`
	XMLNS   string   `xml:"xmlns,attr"`
	Blocks  []Block  `xml:"block"`
}

type Block struct {
	XMLName    xml.Name    `xml:"block"`
	Type       string      `xml:"type,attr"`
	Mutation   *Mutation   `xml:"mutation,omitempty"`
	Fields     []Field     `xml:"field"`
	Values     []Value     `xml:"value"`
	Statements []Statement `xml:"statement"`
	Next       *Next       `xml:"next"`
}

type Field struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type Value struct {
	XMLName xml.Name `xml:"value"`
	Name    string   `xml:"name,attr"`
	Block   Block    `xml:"block"`
}

type Mutation struct {
	XMLName       xml.Name    `xml:"mutation"`
	ItemCount     int         `xml:"items,attr,omitempty"`
	ElseIfCount   int         `xml:"elseif,attr,omitempty"`
	ElseCount     int         `xml:"else,attr,omitempty"`
	LocalNames    []LocalName `xml:"localname"`
	Args          []Arg       `xml:"arg"`
	Key           string      `xml:"key,attr,omitempty"`
	SetOrGet      string      `xml:"set_or_get,attr,omitempty"`
	PropertyName  string      `xml:"property_name,attr,omitempty"`
	IsGeneric     bool        `xml:"is_generic,attr,omitempty"`
	ComponentType string      `xml:"component_type,attr,omitempty"`
	InstanceName  string      `xml:"instance_name,attr,omitempty"`
	EventName     string      `xml:"event_name,attr,omitempty"`
	MethodName    string      `xml:"method_name,attr,omitempty"`
	Mode          string      `xml:"mode,attr,omitempty"`
	Cofounder     string      `xml:"confounder,attr,omitempty"`
	Inline        bool        `xml:"inline,attr,omitempty"`
	Name          string      `xml:"name,attr,omitempty"`
}

type LocalName struct {
	XMLName xml.Name `xml:"localname"`
	Name    string   `xml:"name,attr"`
}

type Statement struct {
	XMLName xml.Name `xml:"statement"`
	Name    string   `xml:"name,attr"`
	Block   *Block   `xml:"block"`
}

type Next struct {
	XMLName xml.Name `xml:"next"`
	Block   *Block   `xml:"block"`
}

type Arg struct {
	Name string `xml:"name,attr"`
}

type Expr interface {
	String() string
	Blockly() Block
	Continuous() bool
	Consumable() bool
}

func (b *Block) String() string {
	return "<" + b.Type + ">"
}

func (b *Block) GetType() string {
	return b.Type
}

func (b *Block) Order() int {
	return 100
}

func (b *Block) SingleValue() Block {
	return b.Values[0].Block
}

func (b *Block) SingleField() string {
	return b.Fields[0].Value
}

func (b *Block) SingleStatement() Statement {
	return b.Statements[0]
}

func (b *Block) Statement() Statement {
	return b.Statements[0]
}
