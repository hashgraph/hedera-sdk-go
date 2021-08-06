package main

import (
	"go/doc"
    "go/ast"
	// "go/format"
	"go/parser"
	"go/token"
	// "io/ioutil"
	"path/filepath"
	"regexp"
)

var protobufTypeRegex = regexp.MustCompile("ProtobufType: .+")
var protobufAccessorRegex = regexp.MustCompile("ProtobufAccessor: .+")

var dir = filepath.Dir(".")
var path = "./"
var fileSet = token.NewFileSet()

func ParseDir(mode parser.Mode) map[string]*ast.Package {
    pkgs, err := parser.ParseDir(fileSet, path, nil, mode)
	if err != nil {
		panic(err)
	}

    return pkgs
}

func GetAstPackage() *ast.Package {
	return ParseDir(0)["hedera"]
}

func GetDocPackage() *ast.Package {
	return ParseDir(parser.ParseComments)["hedera"]

}
func main() {
    astPkg := GetAstPackage()
    docPkg := GetDocPackage()

	documentation := doc.New(docPkg, path, 0)
    structs := StructsFromFiles(astPkg.Files, documentation)

    structs.WriteToFiles()
}
