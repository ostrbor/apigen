package main

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
	Type        string            `yaml:"type"`
	Description *string           `yaml:"description"`
	Items       *Schema           `yaml:"items"`
	Properties  map[string]Schema `yaml:"properties"`
	Required    []string          `yaml:"required"`
	MinLength   *int              `yaml:"minLength"`
	MaxLength   *int              `yaml:"maxLength"`
	Pattern     *string           `yaml:"pattern"`
	Enum        []string          `yaml:"enum"`
	Format      *string           `yaml:"format"`
	Minimum     *int              `yaml:"minimum"`
	Maximum     *int              `yaml:"maximum"`
	MultipleOf  *int              `yaml:"multipleOf"`
	MinItems    *int              `yaml:"minItems"`
	MaxItems    *int              `yaml:"maxItems"`
	UniqueItems *bool             `yaml:"uniqueItems"`
}

type Response struct {
	Content Content `yaml:"content"`
}
