package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

const handlerTmpl = `
package handlers

import (
    "database/sql"
    "fmt"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/arindam923/api-generator/db"
    "github.com/arindam923/api-generator/models"
)

// Get{{.Name}} retrieves a {{.Name}} by ID
func Get{{.Name}}(c *gin.Context) {
    id := c.Param("id")
    {{.VarName}} := &models.{{.Name}}{}
    err := db.Get(db.DB, {{.VarName}}, "SELECT * FROM {{.TableName}} WHERE id = $1", id)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "{{.Name}} not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, {{.VarName}})
}

// List{{.Name}}s retrieves all {{.Name}}s
func List{{.Name}}s(c *gin.Context) {
    var {{.VarNamePlural}} []*models.{{.Name}}
    err := db.Select(db.DB, &{{.VarNamePlural}}, "SELECT * FROM {{.TableName}}")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, {{.VarNamePlural}})
}

// Create{{.Name}} creates a new {{.Name}}
func Create{{.Name}}(c *gin.Context) {
    var {{.VarName}} models.{{.Name}}
    if err := c.ShouldBindJSON(&{{.VarName}}); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    err := db.Insert(db.DB, &{{.VarName}}, "INSERT INTO {{.TableName}} ({{.ColumnsList}}) VALUES ({{.PlaceholdersList}}) RETURNING id", {{.PlaceholderValues}})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, {{.VarName}})
}

// Update{{.Name}} updates an existing {{.Name}}
func Update{{.Name}}(c *gin.Context) {
    id := c.Param("id")
    var {{.VarName}} models.{{.Name}}
    if err := c.ShouldBindJSON(&{{.VarName}}); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    {{.VarName}}.ID, _ = strconv.ParseUint(id, 10, 64)
    _, err := db.Update(db.DB, &{{.VarName}}, "UPDATE {{.TableName}} SET {{.UpdateSetPlaceholders}} WHERE id = $1", {{.UpdatePlaceholderValues}})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, {{.VarName}})
}

// Delete{{.Name}} deletes a {{.Name}} by ID
func Delete{{.Name}}(c *gin.Context) {
    id := c.Param("id")
    _, err := db.Exec(db.DB, "DELETE FROM {{.TableName}} WHERE id = $1", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "{{.Name}} deleted successfully"})
}
`

const routerTmpl = `
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/arindam923/api-generator/handlers"
)

func main() {
    r := gin.Default()

    // {{.Name}} routes
    {{.VarName}}s := r.Group("/{{.PathPlural}}")
    {
        {{.VarName}}s.GET("/", handlers.List{{.Name}}s)
        {{.VarName}}s.GET("/:id", handlers.Get{{.Name}})
        {{.VarName}}s.POST("/", handlers.Create{{.Name}})
        {{.VarName}}s.PUT("/:id", handlers.Update{{.Name}})
        {{.VarName}}s.DELETE("/:id", handlers.Delete{{.Name}})
    }

    r.Run()
}
`

// Schema represents a data model schema
type Schema struct {
	Name      string
	TableName string
	Fields    []*Field
}

// Field represents a field in a data model schema
type Field struct {
	Name string
	Type reflect.Type
	Tag  string
}

func GenerateRESTAPI(schemas []*Schema, outputDir string) error {
	handlerDir := filepath.Join(outputDir, "handlers")
	routerDir := outputDir

	if err := os.MkdirAll(handlerDir, 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(routerDir, 0755); err != nil {
		return err
	}

	for _, schema := range schemas {
		handlerPath := filepath.Join(handlerDir, fmt.Sprintf("%s.go", strings.ToLower(schema.Name)))
		routerPath := filepath.Join(routerDir, "main.go")

		if err := generateHandlerFile(schema, handlerPath); err != nil {
			return err
		}

		if err := generateRouterFile(schema, routerPath); err != nil {
			return err
		}
	}

	return nil
}

func generateColumnsList(schema *Schema) string {
	var columns []string
	for _, field := range schema.Fields {
		columns = append(columns, field.Name)
	}
	return strings.Join(columns, ", ")
}

func generatePlaceholdersList(schema *Schema) string {
	var placeholders []string
	for i := range schema.Fields {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
	}
	return strings.Join(placeholders, ", ")
}

func generatePlaceholderValues(schema *Schema) []interface{} {
	var values []interface{}
	for _, field := range schema.Fields {
		values = append(values, field.Name)
	}
	return values
}

func generateUpdateSetPlaceholders(schema *Schema) string {
	var placeholders []string
	for i, field := range schema.Fields {
		placeholders = append(placeholders, fmt.Sprintf("%s = $%d", field.Name, i+1))
	}
	return strings.Join(placeholders, ", ")
}

func generateHandlerFile(schema *Schema, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl, err := template.New("handler").Parse(handlerTmpl)
	if err != nil {
		return err
	}

	data := struct {
		Name                    string
		VarName                 string
		VarNamePlural           string
		TableName               string
		ColumnsList             string
		PlaceholdersList        string
		PlaceholderValues       []interface{}
		UpdateSetPlaceholders   string
		UpdatePlaceholderValues []interface{}
	}{
		Name:                    schema.Name,
		VarName:                 strings.ToLower(schema.Name),
		TableName:               schema.TableName,
		VarNamePlural:           pluralize(strings.ToLower(schema.Name)),
		ColumnsList:             generateColumnsList(schema),
		PlaceholdersList:        generatePlaceholdersList(schema),
		PlaceholderValues:       generatePlaceholderValues(schema),
		UpdateSetPlaceholders:   generateUpdateSetPlaceholders(schema),
		UpdatePlaceholderValues: append(generatePlaceholderValues(schema), uint64(0)),
	}

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("error executing handler template: %v", err)
	}

	return nil
}

func generateRouterFile(schema *Schema, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl, err := template.New("router").Parse(routerTmpl)
	if err != nil {
		return err
	}

	data := struct {
		Name       string
		VarName    string
		PathPlural string
	}{
		Name:       schema.Name,
		VarName:    strings.ToLower(schema.Name),
		PathPlural: strings.ToLower(pluralize(schema.Name)),
	}
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("error executing router template: %v", err)
	}

	return nil
}

func pluralize(str string) string {
	if strings.HasSuffix(str, "y") {
		return strings.TrimSuffix(str, "y") + "ies"
	}
	return str + "s"
}

// ParseSchemas parses the given Go structs and returns their schema representations
func ParseSchemas(models ...interface{}) ([]*Schema, error) {
	var schemas []*Schema

	for _, model := range models {
		schema, err := parseSchema(model)
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, schema)
	}

	return schemas, nil
}

func parseSchema(model interface{}) (*Schema, error) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid model type: %v", modelType.Kind())
	}

	schema := &Schema{
		Name:   modelType.Name(),
		Fields: make([]*Field, 0),
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldType := field.Type

		tag := field.Tag.Get("json")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}

		schema.Fields = append(schema.Fields, &Field{
			Name: field.Name,
			Type: fieldType,
			Tag:  tag,
		})
	}

	return schema, nil
}
