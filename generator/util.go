package main

import (
	"fmt"
	"go/doc"
	"io/ioutil"
	"strings"
	"unicode"
)

func ReadFileToString(fileName string) string {
	s, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	return string(s)
}

func CreateGeneratedFileName(name string) string {
	name = strings.Split(name, ".go")[0]
	return fmt.Sprintf("%s_generated.go", name)
}

func FindProtobufNameInDocumentation(documentation string) *string {
	matches := protobufTypeRegex.FindAllString(documentation, -1)
	if matches == nil || len(matches) != 1 {
		return nil
	}

	return &strings.Split(matches[0], " ")[1]
}

func FindProtobufAccessorInDocumentation(documentation string) *string {
	matches := protobufAccessorRegex.FindAllString(documentation, -1)
	if matches == nil || len(matches) != 1 {
		return nil
	}

	return &strings.Split(matches[0], " ")[1]
}

func FindDocType(documentation *doc.Package, typeName string) string {
	d := ""
	for _, typ := range documentation.Types {
		if typ.Name == typeName {
			d = typ.Doc
		}
	}

	if d == "" {
		panic(fmt.Sprintf("Unable to find documentation for type name %s", typeName))
	}

	return d
}

func LowerInitial(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]

	}

	return ""
}

func UpperInitial(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]

	}

	return ""
}
