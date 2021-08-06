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
    "fmt"
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

    fmt.Printf("Structs: %+v\n", structs)

    structs.WriteToFiles()
	for _, structure := range structs.structs {
		data := []byte(structure.String())
		src, err := format.Source(data)
		if err != nil {
			panic(err)
		}
	
		output := filepath.Join(Dir, structure.fileName)
		err = ioutil.WriteFile(output, src, 0644)
		if err != nil {
			panic(err)
		}
	}
}

