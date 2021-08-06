package main

import (
	"go/doc"
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

func main() {
	pkgs, err := parser.ParseDir(fileSet, path, nil, 0)
	if err != nil {
		panic(err)
	}

	// Used for AST
	pkg := pkgs["hedera"]

	pkgs, err = parser.ParseDir(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// Used for Comments
	// The reason we need this is because `doc.New()` takes over the `pkg` and completely removes the AST
	docPkg := pkgs["hedera"]

	documentation := doc.New(docPkg, "./", 0)
	structs := Structs{
        structs: StructsFromFiles(pkg.Files, documentation),
    }

    fmt.Printf("Structs: %+v\n", structs)
	// for _, structure := range structs.structs {
	// 	data := []byte(structure.String())
	// 	src, err := format.Source(data)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 
	// 	output := filepath.Join(Dir, structure.fileName)
	// 	err = ioutil.WriteFile(output, src, 0644)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
}

