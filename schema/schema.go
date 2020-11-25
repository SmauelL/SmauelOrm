package schema

import (
	"go/ast"
	"reflect"
	"smauelOrm/dialect"
)

// Filed represents a column of database
type Filed struct {
	Name string
	Type string
	Tag  string
}

// Schema represents a table of database
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Filed
	FiledNames []string
	filedMap   map[string]*Filed
}

func (schema *Schema) GetFiled(name string) *Filed {
	return schema.filedMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		filedMap: make(map[string]*Filed),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			filed := &Filed{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("smauelOrm"); ok {
				filed.Tag = v
			}
			schema.Fields = append(schema.Fields, filed)
			schema.FiledNames = append(schema.FiledNames, p.Name)
			schema.filedMap[p.Name] = filed
		}
	}
	return schema
}

func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
