package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type Field struct {
	name   string
	ty     FieldType
	config FieldConfig
}

type FieldConfig struct {
	name         string
	getter       bool
	setter       bool
	singular     bool
	toProtobuf   bool
	fromProtobuf bool
	protobufName string
}

type FieldType struct {
	reference   bool
    array       bool
	packageName string
	name        string
}

func FieldConfigFromTag(fieldTag *ast.BasicLit) FieldConfig {
	fieldConfig := FieldConfig{}
	if fieldTag == nil || fieldTag.Kind != token.STRING {
		return fieldConfig
	}

	tags := strings.Split(fieldTag.Value, " ")

	if len(tags) == 0 {
		return fieldConfig
	}

	var options []string
	for _, t := range tags {
		if strings.HasPrefix(t, "`hedera:") {
			options = strings.Split(strings.Split(t, "\"")[1], ",")
			break
		}
	}

	if options == nil {
		return FieldConfig{}
	}

	for _, option := range options {
		o := strings.Split(option, "=")
		switch o[0] {
		case "setter":
			fieldConfig.setter = true
		case "getter":
			fieldConfig.getter = true
		case "singular":
			fieldConfig.singular = true
		case "toProtobuf":
			fieldConfig.toProtobuf = true
		case "fromProtobuf":
			fieldConfig.fromProtobuf = true
		case "protobufName":
			if len(o) != 2 {
				panic("protobufName used, but no value provided")
			}

			fieldConfig.protobufName = o[1]
		}
	}

	return fieldConfig
}

func FieldsFromExpr(expr ast.Expr) []Field {
	switch expr.(type) {
	case *ast.StructType:
        return FieldsFromStructType(expr.(*ast.StructType))
    default:
        return make([]Field, 0)
	}
}

func FieldsFromStructType(expr *ast.StructType) []Field {
    fields := make([]Field, 0)

    for _, field := range expr.Fields.List {
        f := FieldFromField(field)
        if f == nil {
            return nil
        }

        fields = append(fields, *f)
    }

    return fields
}

func FieldFromField(field *ast.Field) *Field {
    config := FieldConfigFromTag(field.Tag)
    ty := FieldTypeFromType(field.Type)

    if len(field.Names) != 1 {
        return nil
    }

    return &Field{
        name:   field.Names[0].Name,
        ty:     ty,
        config: config,
    }
}

func FieldTypeFromIdent(ident *ast.Ident) FieldType {
	return FieldType{
		name: ident.Name,
	}
}

func FieldTypeFromStarExpr(expr *ast.StarExpr) FieldType {
	switch expr.X.(type) {
	case *ast.SelectorExpr:
		fieldType := FieldTypeFromSelectorExpr(expr.X.(*ast.SelectorExpr))
		fieldType.reference = true
		return fieldType
	default:
		panic(fmt.Sprintf("Unhandled field type: %T", expr.X))
	}
}

func FieldTypeFromSelectorExpr(expr *ast.SelectorExpr) FieldType {
	switch expr.X.(type) {
	case *ast.Ident:
		return FieldType{
			packageName: fmt.Sprintf("%s", expr.X),
			name:        fmt.Sprintf("%s", expr.Sel),
		}
	default:
		panic(fmt.Sprintf("Unhandled field type: %T", expr.X))
	}
}

func FieldTypeFromArrayType(expr *ast.ArrayType) FieldType {
	return FieldType{
        array:     true,
		name:      fmt.Sprintf("%s", expr.Elt),
	}
}

func FieldTypeFromType(expr ast.Expr) FieldType {
	switch expr.(type) {
	case *ast.Ident:
		return FieldTypeFromIdent(expr.(*ast.Ident))
	case *ast.StarExpr:
		return FieldTypeFromStarExpr(expr.(*ast.StarExpr))
	case *ast.SelectorExpr:
		return FieldTypeFromSelectorExpr(expr.(*ast.SelectorExpr))
	case *ast.ArrayType:
		return FieldTypeFromArrayType(expr.(*ast.ArrayType))
	default:
		panic(fmt.Sprintf("Unhandled field type: %T", expr))
	}
}

func (field Field) Replacer() *strings.Replacer {
    name := field.name
    if field.config.singular && strings.HasSuffix(name, "s") {
        name = name[0:len(name)-1]
    }

    ty := field.ty.String(field.config)

    return strings.NewReplacer("<field.name>", name, "<field.type>", ty)
}

func (ty FieldType) String(config FieldConfig) string {
    s := ""

    if ty.reference {
        s += "*"
    }

    if ty.array && !config.singular {
        s += "[]"
    }

    if ty.packageName != "" {
        s += ty.packageName
    }

    if config.singular && strings.HasSuffix(ty.name, "s") {
        s += ty.name[0:len(ty.name)-1]
    } else {
        s += ty.name
    }

    return s
}
