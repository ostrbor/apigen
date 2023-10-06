package main

import (
	"io"
	"os"
	"strings"
	"text/template"
)

type goStruct struct {
	Name   string
	Fields []goField
}

type goField struct {
	Name   string
	Typ    string
	JSName string
}

// func toGoStructs converts an OpenAPI Schema into Go structs.
func toGoStruct(name string, schema Schema) (res []goStruct) {
	if schema.Type != "object" {
		panic("toGoStruct only supports object schemas")
	}
	if name == "" {
		name = "Autogenerated"
	}
	st := goStruct{Name: name}
	for propName, propSchema := range schema.Properties {
		if propSchema.Type == "object" {
			st.Fields = append(st.Fields, goField{
				Name:   snakeToCamel(propName),
				Typ:    snakeToCamel(propName),
				JSName: propName,
			})
			res = append(res, toGoStruct(snakeToCamel(propName), propSchema)...)
		} else {
			st.Fields = append(st.Fields, goField{
				Name:   snakeToCamel(propName),
				Typ:    toGoType(propSchema),
				JSName: propName,
			})
		}
	}
	res = append(res, st)
	return res
}

func toGoType(schema Schema) (goType string) {
	switch schema.Type {
	case "string", "number", "integer", "boolean":
		goType = goTypes[schema.Type]
	case "array":
		itemTyp := (*schema.Items).Type
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

func generateCode(s goStruct) string {
	f, err := os.Open("struct.tmpl")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	text, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	tmpl := template.New("GoStruct")
	tmpl, err = tmpl.Parse(string(text))
	if err != nil {
		panic(err)
	}
	buf := new(strings.Builder)
	err = tmpl.Execute(buf, s)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
