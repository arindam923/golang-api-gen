package main

import (
	"fmt"

	"github.com/arindam923/api-generator/generator"
	"github.com/arindam923/api-generator/models"
)

func main() {
	schemas, err := generator.ParseSchemas(models.User{}, models.Post{})
	if err != nil {
		fmt.Println(err)
		return
	}

	outputDir := "api/rest"
	if err := generator.GenerateRESTAPI(schemas, outputDir); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("REST API generated successfully!")
	for _, schema := range schemas {
		fmt.Printf("Schema: %s\n", schema.Name)
		for _, field := range schema.Fields {
			fmt.Printf("\tField: %s (%s) `%s`\n", field.Name, field.Type, field.Tag)
		}
	}
}
