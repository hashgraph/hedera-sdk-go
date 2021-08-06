package main

import (
	"go/ast"
	"go/doc"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Structs struct {
	structs []Struct
}

type Struct struct {
	name          string
	fileName      string
	fields        []Field
	protoName     string
	protoAccessor string
}

var request string

func init() {
	s, err := ioutil.ReadFile("./generator/templates/request.txt")
	if err != nil {
		panic(err)
	}

	request = string(s)
}

func StructsFromFiles(files map[string]*ast.File, documentation *doc.Package) Structs {
	structs := Structs{
		structs: make([]Struct, 0),
	}

	for fileName, file := range files {
		structs.structs = append(structs.structs, StructsFromFile(fileName, file, documentation)...)
	}

	return structs
}

func StructsFromFile(fileName string, file *ast.File, documentation *doc.Package) []Struct {
	structs := make([]Struct, 0)

	for _, decl := range file.Decls {
		structs = append(structs, StructsFromDecl(fileName, decl, documentation)...)
	}

	return structs
}

func StructsFromDecl(fileName string, decl ast.Decl, documentation *doc.Package) []Struct {
	switch decl.(type) {
	case *ast.GenDecl:
		return StructsFromGenDecl(fileName, decl.(*ast.GenDecl), documentation)
	default:
		return make([]Struct, 0)
	}
}

func StructsFromGenDecl(fileName string, decl *ast.GenDecl, documentation *doc.Package) []Struct {
	structs := make([]Struct, 0)

	for _, spec := range decl.Specs {
		structure := StructFromSpec(fileName, spec, documentation)
		if structure == nil {
			continue
		}

		structs = append(structs, *structure)
	}

	return structs
}

func StructFromSpec(fileName string, spec ast.Spec, documentation *doc.Package) *Struct {
	switch spec.(type) {
	case *ast.TypeSpec:
		return StructFromTypeSpec(fileName, spec.(*ast.TypeSpec), documentation)
	default:
		return nil
	}
}

func StructFromTypeSpec(fileName string, spec *ast.TypeSpec, documentation *doc.Package) *Struct {
	if !((strings.HasSuffix(spec.Name.Name, "Transaction") && len(spec.Name.Name) > 12) || (strings.HasSuffix(spec.Name.Name, "Query") && len(spec.Name.Name) > 6)) {
		return nil
	}

	if spec.Name.Name != "AccountCreateTransaction" {
		return nil
	}

	d := FindDocType(documentation, spec.Name.Name)

	protoName := FindProtobufNameInDocumentation(d)
	protoAccessor := FindProtobufAccessorInDocumentation(d)
	fileName = CreateGeneratedFileName(fileName)

	if protoName == nil || protoAccessor == nil {
		return nil
	}

	return &Struct{
		name:          spec.Name.Name,
		fileName:      fileName,
		fields:        FieldsFromExpr(spec.Type),
		protoName:     *protoName,
		protoAccessor: *protoAccessor,
	}
}

func (structs Structs) String() string {
	s := ""

	for _, structure := range structs.structs {
		s += structure.String()
	}

	return s
}

func (structure Struct) String() string {
	replacer := strings.NewReplacer(
		"<imports>", GenerateImports(structure),
		"<setters>", GenerateSetters(structure),
		"<getters>", GenerateGetters(structure),
		"<fromProtobuf>", GenerateFromProtobufs(structure),
		"<toProtobuf>", GenerateToProtobufs(structure),
		"<validateChecksums>", GenerateValidateChecksums(structure),
	)

	if replacer == nil {
		return ""
	}

	s := replacer.Replace(request)

	if formatted, err := format.Source([]byte(s)); err == nil {
		s = string(formatted)
	}

	return s
}

func (structs Structs) WriteToFiles() {
	for _, structure := range structs.structs {
		data := []byte(structure.String())
		src, err := format.Source(data)
		if err != nil {
			panic(err)
		}

		output := filepath.Join(dir, structure.fileName)
		err = ioutil.WriteFile(output, src, 0644)
		if err != nil {
			panic(err)
		}
	}
}
