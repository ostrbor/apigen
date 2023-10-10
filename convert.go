package main

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

func objToStruct(name string, schema Schema) (res []string) {
	if schema.Type != "object" {
		return nil
	}
	if name == "" {
		name = "Autogenerated"
	}

	buf := new(strings.Builder)
	w := tabwriter.NewWriter(buf, 0, 0, 0, ' ', 0)

	fmt.Fprintf(w, "type %s struct {\n", name)
	for propName, propSchema := range schema.Properties {
		switch propSchema.Type {
		case "object":
			field := toCamel(propName)
			jsonTag := toJSONTag(propName)
			fmt.Fprintf(w, "    %s \t%s \t%s\n", field, field, jsonTag)
			res = append(res, objToStruct(field, propSchema)...)
		case "array":
			if (*propSchema.Items).Type == "object" {
				field := toCamel(propName)
				objName := toSingle(field)
				jsonTag := toJSONTag(propName)
				fmt.Fprintf(w, "    %s \t[]%s \t%s\n", field, objName, jsonTag)
				res = append(res, objToStruct(objName, *propSchema.Items)...)
			} else {
				field := toCamel(propName)
				typ := goTypes[(*propSchema.Items).Type]
				jsonTag := toJSONTag(propName)
				fmt.Fprintf(w, "    %s \t[]%s \t%s\n", field, typ, jsonTag)
			}
		default:
			field := toCamel(propName)
			typ := goTypes[propSchema.Type]
			jsonTag := toJSONTag(propName)
			fmt.Fprintf(w, "    %s \t%s \t%s\n", field, typ, jsonTag)
		}
	}
	fmt.Fprintf(w, "}\n")

	w.Flush()
	return append(res, buf.String())
}

var goTypes = map[string]string{
	"string":  "string",
	"number":  "int",
	"integer": "int",
	"boolean": "bool",
}

func toSingle(s string) string {
	return strings.TrimSuffix(s, "s")
}

func toCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func toJSONTag(name string) string {
	return fmt.Sprintf("`json:\"%s\"`", name)
}
