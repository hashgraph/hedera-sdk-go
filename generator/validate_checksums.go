package main

import (
	"io/ioutil"
	"strings"
)

var validateChecksumsFunction string
var validateChecksumsId string

func init() {
	s, err := ioutil.ReadFile("./generator/templates/validate_checksums/function.txt")
	if err != nil {
		panic(err)
	}

	validateChecksumsFunction = string(s)

	s, err = ioutil.ReadFile("./generator/templates/validate_checksums/id.txt")
	if err != nil {
		panic(err)
	}

	validateChecksumsId = string(s)
}

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
