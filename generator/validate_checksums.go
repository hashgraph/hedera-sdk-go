package main

import (
	"strings"
)

var validateChecksumsFunction string = ReadFileToString("./generator/templates/validate_checksums/function.txt")
var validateChecksumsId string = ReadFileToString("./generator/templates/validate_checksums/id.txt")

func GenerateValidateChecksums(structure Struct) string {
	s := ""

	for _, field := range structure.fields {
		if field.ty.array || !strings.HasSuffix(field.ty.name, "ID") {
			continue
		}

		replacer := field.Replacer()

		if replacer == nil {
			continue
		}

		s += replacer.Replace(validateChecksumsId)
	}

	replacer := strings.NewReplacer(
		"<this.type>", structure.name,
		"<conditionals>", s,
	)

	return replacer.Replace(validateChecksumsFunction)
}
