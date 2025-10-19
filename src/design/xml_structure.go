package design

import "encoding/xml"

type XmlRoot struct {
	XMLName xml.Name  `xml:"xml"`
	XMLNS   string    `xml:"xmlns,attr"`
	Screen  Component `xml:"Screen"`
}

type Component struct {
	XMLName    xml.Name          `xml:""`
	Id         string            `xml:"id,attr,omitempty"`
	Type       string            `xml:""`
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
