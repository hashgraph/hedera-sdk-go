package main

import (
    "io/ioutil"
    "strings"
)

var getter string
var getterSingular string

func init() {
    s, err := ioutil.ReadFile("./generator/templates/getter.txt")
    if err != nil {
        panic(err)
    }

    getter = string(s)

    s, err = ioutil.ReadFile("./generator/templates/getter_singular.txt")
    if err != nil {
        panic(err)
    }

    getterSingular = string(s)
}

func GenerateGetters(structure Struct) string {
    s := ""

    for _, field := range structure.fields {
        if !field.config.getter {
            continue
        }

        replacer := field.Replacer()

        if replacer == nil {
            continue
        }

        if field.config.singular {
            s += strings.ReplaceAll(replacer.Replace(getterSingular), "<this.type>", structure.name)
        } else {
            s += strings.ReplaceAll(replacer.Replace(getter), "<this.type>", structure.name)
        }
    }

    return s
}
