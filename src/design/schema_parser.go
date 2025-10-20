package design

import (
	"encoding/json"
	"encoding/xml"
)

type SchemaParser struct {
	schemaJson string
}

func NewSchemaParser(schemaJson string) *SchemaParser {
	return &SchemaParser{schemaJson: schemaJson}
}

func (p *SchemaParser) ConvertSchemaToXml() (string, error) {
	var jsonStruct map[string]interface{}
	err := json.Unmarshal([]byte(p.schemaJson), &jsonStruct)
	if err != nil {
		return "", err
	}
	properties := jsonStruct["Properties"].(map[string]interface{})
	screenId := properties["$Name"].(string)

	schemaComponents := properties["$Components"].([]interface{})
	var xmlChildren []Component
	for _, schemaComponent := range schemaComponents {
		xmlChildren = append(xmlChildren, schemaComponentToXml(schemaComponent.(map[string]interface{})))
	}

	root := Component{
		XMLName:    xml.Name{Local: "Screen"},
		Id:         screenId,
		Type:       "Form",
		Properties: filterDesignerProperties(properties),
		Children:   xmlChildren,
	}
	result, err := xml.MarshalIndent(root, "", "  ")
	return string(result), err
}

func schemaComponentToXml(schemaJson map[string]interface{}) Component {
	compType := schemaJson["$Type"].(string)
	schemaComponents := schemaJson["$Components"]
	var xmlChildren []Component
	if schemaComponents != nil {
		for _, schemaComponent := range schemaComponents.([]interface{}) {
			xmlChildren = append(xmlChildren, schemaComponentToXml(schemaComponent.(map[string]interface{})))
		}
	}
	if len(xmlChildren) == 0 {
		xmlChildren = nil
	}

	return Component{
		XMLName:    xml.Name{Local: compType},
		Id:         schemaJson["$Name"].(string),
		Type:       compType,
		Properties: filterDesignerProperties(schemaJson),
		Children:   xmlChildren,
	}
}

func filterDesignerProperties(componentProps map[string]interface{}) map[string]string {
	filteredProperties := make(map[string]string)
	for key, value := range componentProps {
		if len(key) > 0 && key[0] == '$' {
			continue
		}
		filteredProperties[key] = value.(string)
	}
	return filteredProperties
}
