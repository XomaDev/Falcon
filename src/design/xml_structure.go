package design

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
)

type XmlRoot struct {
	XMLName xml.Name  `xml:"xml"`
	XMLNS   string    `xml:"xmlns,attr"`
	Screen  Component `xml:"Screen"`
}

type Component struct {
	XMLName    xml.Name          `xml:""`
	Id         string            `xml:"id,attr,omitempty"`
	Type       string            `xml:"-"`
	Properties map[string]string `xml:"-"`
	Children   []Component       `xml:",any"`
}

func (c *Component) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	c.XMLName = start.Name
	c.Type = start.Name.Local
	c.Properties = make(map[string]string)

	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			c.Id = attr.Value
		} else {
			c.Properties[attr.Name.Local] = attr.Value
		}
	}
	for {
		tok, err := d.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			var child Component
			if err := d.DecodeElement(&child, &tok); err != nil {
				return err
			}
			c.Children = append(c.Children, child)
		case xml.EndElement:
			if tok.Name.Local == start.Name.Local {
				return nil
			}
		}
	}
	return nil
}

// WriteXML manually converts Component structure to XML, a workaround for now, since
// Go lang does not support self-closing tags
func (c *Component) WriteXML(w io.Writer, indent int) error {
	indentStr := strings.Repeat("  ", indent)

	// Start tag
	tag := indentStr + "<" + c.Type
	if c.Id != "" {
		var buf bytes.Buffer
		err := xml.EscapeText(&buf, []byte(c.Id))
		if err != nil {
			return err
		}
		tag += ` id="` + buf.String() + `"`
	}

	for k, v := range c.Properties {
		if k == "id" || k == "type" {
			continue
		}
		var buf bytes.Buffer
		err := xml.EscapeText(&buf, []byte(v))
		if err != nil {
			return err
		}
		tag += ` ` + k + `="` + buf.String() + `"`
	}

	if len(c.Children) == 0 {
		tag += "/>\n"
		_, err := w.Write([]byte(tag))
		return err
	}

	// Open tag
	tag += ">\n"
	if _, err := w.Write([]byte(tag)); err != nil {
		return err
	}

	// Recursively write children
	for _, child := range c.Children {
		if err := child.WriteXML(w, indent+1); err != nil {
			return err
		}
	}

	// Close tag
	_, err := w.Write([]byte(indentStr + "</" + c.Type + ">\n"))
	return err
}
