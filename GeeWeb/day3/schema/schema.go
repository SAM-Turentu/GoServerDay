package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

/*
object <==> table
*/

type Field struct {
	Name string // 字段名
	Type string // 字段类型
	Tag  string // 约束条件
}

// Schema 映射对象 Model，表名 Name， 字段 Fields
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(), // 获取结构体的名称作为表名
		fieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,                                              // 字段名
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))), // 字段类型
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v // 约束条件
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
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
