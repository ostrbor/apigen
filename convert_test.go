package main

import (
	"gopkg.in/yaml.v2"
	"testing"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func TestConvert(t *testing.T) {
	yamlSchema := `
type: object
properties:
  str:
   type: string
  int:
    type: integer
  bool:
    type: boolean
  arr:
    type: array
    items:
      type: string
`
	var schema Schema
	err := yaml.Unmarshal([]byte(yamlSchema), &schema)
	check(err)
	s := toGoStruct("", schema)
	if len(s.fields) != 4 {
		t.Errorf("Expected 4 fields, got %d", len(s.fields))
	}
}
