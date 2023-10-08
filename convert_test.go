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

func Test_objToStruct(t *testing.T) {
	yamlSchema := `
type: object
properties:
  name:
   type: string
  age:
    type: integer
  is_active:
    type: boolean
  phones:
    type: array
    items:
      type: string
  address:
    type: object
    properties:
      street:
        type: string
`
	var schema Schema
	err := yaml.Unmarshal([]byte(yamlSchema), &schema)
	check(err)
	ss := objToStruct("", schema)
	for _, s := range ss {
		println(s)
	}
}
