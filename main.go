package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var (
	openapiPath = flag.String("path", "", "Set the OpenAPI path to process")
	//openapiVerb = flag.String("verb", "", "Set the OpenAPI verb to process")
	openapiFile = flag.String("file", "openapi.yaml", "Specify the OpenAPI file")
)

func main() {
	flag.Parse()
	if _, err := os.Stat(*openapiFile); os.IsNotExist(err) {
		log.Printf("file '%s' does not exist\n", *openapiFile)
		os.Exit(1)
	}

	f, err := os.ReadFile(*openapiFile)
	if err != nil {
		log.Printf("error reading file '%s': %v\n", *openapiFile, err)
		os.Exit(1)
	}

	api := OpenAPI{}
	if err := yaml.Unmarshal(f, &api); err != nil {
		log.Printf("error unmarshalling file '%s': %v\n", *openapiFile, err)
		os.Exit(1)
	}

	for path, item := range api.Paths {
		if *openapiPath != "" && *openapiPath != path {
			continue
		}
		if item.Get != nil {
			req := item.Get.RequestBody.Content.ApplicationJSON.Schema
			for _, s := range objToStruct("GetRequest", req) {
				println(s)
			}
			for _, resp := range item.Get.Responses {
				for _, s := range objToStruct("GetResponse", resp.Content.ApplicationJSON.Schema) {
					println(s)
				}
			}
		}
		if item.Post != nil {
			req := item.Post.RequestBody.Content.ApplicationJSON.Schema
			for _, s := range objToStruct("PostRequest", req) {
				println(s)
			}
			for _, resp := range item.Post.Responses {
				for _, s := range objToStruct("PostResponse", resp.Content.ApplicationJSON.Schema) {
					println(s)
				}
			}
		}
	}
}
