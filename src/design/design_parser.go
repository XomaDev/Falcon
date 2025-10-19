package design

import (
	"encoding/json"
	"encoding/xml"
)

func ConvertXmlToSchema(xmlContent string) (string, error) {
	var root XmlRoot
	if err := xml.Unmarshal([]byte(xmlContent), &root); err != nil {
		panic(err)
	}
	var components []interface{}
	for _, child := range root.Screen.Children {
		components = append(components, componentToJson(child))
	}

	screen := root.Screen
	props := map[string]interface{}{
		"$Name":       screen.Id,
		"$Type":       "Form",
		"$Version":    "31",
		"$Components": components,
	}
	// add Screen's properties here
	for k, v := range screen.Properties {
		props[k] = v
	}

	schema := map[string]interface{}{
		"authURL":    []interface{}{"ai2.appinventor.mit.edu"},
		"YaVersion":  "200",
		"Source":     "Form",
		"Properties": props,
	}
	jsonBytes, _ := json.MarshalIndent(schema, "", "  ")
	return string(jsonBytes), nil
}

func componentToJson(component Component) interface{} {
	var children []interface{}
	for _, child := range component.Children {
		children = append(children, componentToJson(child))
	}
	schema := map[string]interface{}{
		"$Name":    component.Id,
		"$Type":    component.Type,
		"$Version": "0",
	}
	if len(children) > 0 {
		schema["$Components"] = children
	}
	for k, v := range component.Properties {
		schema[k] = v
	}
	return schema
}
