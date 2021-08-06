package main

import (
	"fmt"
)

var imports map[string]string

func init() {
	imports = make(map[string]string)

	imports["time"] = `"time"`
	imports["proto"] = `"github.com/hashgraph/hedera-sdk-go/v2/proto"`
}

func GenerateImports(structure Struct) string {
	s := ""
	imports_ := make(map[string]bool, 0)

	for _, field := range structure.fields {
		if field.ty.packageName == "" {
			continue
		}

		v, ok := imports[field.ty.packageName]

		if !ok {
			panic(fmt.Sprintf("found field which uses package %s, but could not resolve which package should be imported", field.ty.packageName))
		}

		imports_[v] = true
	}

	for import_, _ := range imports_ {
		s += import_ + "\n"
	}

	return s
}
