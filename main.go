package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
)

type OpenAPI struct {
	Paths map[string]PathItem `yaml:"paths"`
}

type PathItem struct {
	Get  *Operation `yaml:"get"`
	Post *Operation `yaml:"post"`
}

type Operation struct {
	Summary     string              `yaml:"summary"`
	RequestBody RequestBody         `yaml:"requestBody"`
	Responses   map[string]Response `yaml:"responses"`
}

type RequestBody struct {
	Content `yaml:"content"`
}

type Content struct {
	ApplicationJSON `yaml:"application/json"`
}

type ApplicationJSON struct {
	Schema `yaml:"schema"`
}

type Schema struct {
	Type       string              `yaml:"type"`
	Properties map[string]Property `yaml:"properties"`
}

type Property struct {
	Type         string              `yaml:"type"`
	Description  string              `yaml:"description"`
	MinLength    *int                `yaml:"minLength,omitempty"`
	MaxLength    *int                `yaml:"maxLength,omitempty"`
	Pattern      string              `yaml:"pattern,omitempty"`
	Enum         []string            `yaml:"enum,omitempty"`
	Format       string              `yaml:"format,omitempty"`
	Items        *Schema             `yaml:"items,omitempty"`
	NestedObject map[string]Property `yaml:"properties,omitempty"`
}

type Response struct {
	Description string            `yaml:"description"`
	Content     Content           `yaml:"content"`
	Headers     map[string]Header `yaml:"headers"`
}

type Header struct {
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
}

func main() {
	yamlData, err := os.ReadFile("openapi.yaml")
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	var openAPI OpenAPI
	err = yaml.Unmarshal(yamlData, &openAPI)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	for path, items := range openAPI.Paths {
		if items.Post != nil {
			fmt.Printf("\n===POST %s: %s===\n\n", path, items.Post.Summary)
			fmt.Println("// REQUEST")
			printGoType(pathToCamel(path)+"Request", items.Post.RequestBody.Content.ApplicationJSON.Schema)
			for code, resp := range items.Post.Responses {
				if strings.HasPrefix(code, "2") {
					fmt.Printf("// RESPONSE %s\n", code)
					printGoType(pathToCamel(path)+"Response", resp.Content.ApplicationJSON.Schema)
					fmt.Println()
				}
			}
		}
		if items.Get != nil {
			fmt.Printf("\n===GET: %s===\n\n", items.Get.Summary)
			for code, resp := range items.Get.Responses {
				if strings.HasPrefix(code, "2") {
					fmt.Printf("// RESPONSE %s\n", code)
					printGoType(pathToCamel(path)+"Response", resp.Content.ApplicationJSON.Schema)
					fmt.Println()
				}
			}
		}
	}
}

func printGoType(name string, s Schema) {
	switch s.Type {
	case "string", "number", "integer", "boolean":
		fmt.Printf("type %s %s\n", name, goTypes[s.Type])
	case "array":
		fmt.Printf("type %s []%s\n", name, goTypes[s.Type])
	case "object":
		printGoStruct(name, s.Properties)
	}
}

func printGoStruct(structName string, props map[string]Property) {
	fmt.Printf("type %s struct {\n", structName)
	nested := make(map[string]map[string]Property, 0)
	items := make(map[string]Schema, 0)
	for name, prop := range props {
		var fieldType string
		var fieldName = snakeToCamel(name)
		switch prop.Type {
		case "object":
		case "array":
			fieldType = "[]" + strings.TrimSuffix(fieldName, "s")
		default:
			fieldType = goTypes[prop.Type]
		}
		fmt.Printf("\t%s %s `json:\"%s\"`\n", fieldName, fieldType, name)
		if prop.NestedObject != nil {
			nested[fieldName] = prop.NestedObject
		}
		if prop.Items != nil {
			items[strings.TrimSuffix(fieldName, "s")] = *prop.Items
		}
	}
	if structName != "" {
		fmt.Println("}\n")
	}
	for structName, nestedProps := range nested {
		printGoStruct(structName, nestedProps)
	}
	for structName, schema := range items {
		printGoType(structName, schema)
		fmt.Println()
	}
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

func firstSegment(path string) string {
	segments := strings.Split(path, "/")
	if len(segments) > 1 {
		return segments[1]
	}
	return ""
}

func pathToCamel(path string) string {
	s := firstSegment(path)
	parts := strings.Split(s, "-")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}
