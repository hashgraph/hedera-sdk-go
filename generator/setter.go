package main

import (
    "io/ioutil"
    "strings"
)

var setter string
var setterSingular string

func init() {
    s, err := ioutil.ReadFile("./generator/templates/setter.txt")
    if err != nil {
        panic(err)
    }

    setter = string(s)

    s, err = ioutil.ReadFile("./generator/templates/setter_singular.txt")
    if err != nil {
        panic(err)
    }

    setterSingular = string(s)
}

func GenerateSetters(structure Struct) string {
    s := ""

    for _, field := range structure.fields {
        if !field.config.setter {
            continue
        }

        replacer := field.Replacer()

        if field.config.singular {
            s += strings.ReplaceAll(replacer.Replace(setterSingular), "<this.type>", structure.name)
        } else {
            s += strings.ReplaceAll(replacer.Replace(setter), "<this.type>", structure.name)
        }
    }

    return s
}
