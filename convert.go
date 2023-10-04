package main

import "strings"

type goStruct struct {
	name   string
	fields []goField
}

type goField struct {
	name    string
	typ     string
	jsonTag string
}

// func toGoStructs converts an OpenAPI Schema into Go structs.
func toGoStruct(name string, schema Schema) (res goStruct) {
	if name == "" {
		name = "Schema"
	}
	res.name = name
	for propName, prop := range schema.Properties {
		res.fields = append(res.fields, goField{
			name:    snakeToCamel(propName),
			typ:     toGoType(prop),
			jsonTag: propName,
		})
	}
	return res
}

// toGoType converts an OpenAPI Property into a Go field type and struct jsonTag.
func toGoType(property Property) (goType string) {
	switch property.Type {
	case "string", "number", "integer", "boolean":
		goType = goTypes[property.Type]
	case "array":
		itemTyp := (*property.Items).Type
		goType = "[]" + goTypes[itemTyp]
	}
	return goType
}

var goTypes = map[string]string{
	"string":  "string",
	"number":  "int",
	"integer": "int",
	"boolean": "bool",
}

func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}
