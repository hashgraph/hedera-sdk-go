package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var to_protobuf_function string
var to_protobuf_bool string
var to_protobuf_hbar string
var to_protobuf_id string
var to_protobuf_key string
var to_protobuf_time_duration string
var to_protobuf_string string

func init() {
	s, err := ioutil.ReadFile("./generator/templates/to_protobuf/function.txt")
	if err != nil {
		panic(err)
	}

	to_protobuf_function = string(s)

	s, err = ioutil.ReadFile("./generator/templates/to_protobuf/types/bool.txt")
	if err != nil {
		panic(err)
	}

	to_protobuf_bool = string(s)

	s, err = ioutil.ReadFile("./generator/templates/to_protobuf/types/hbar.txt")
	if err != nil {
		panic(err)
	}

	to_protobuf_hbar = string(s)

	s, err = ioutil.ReadFile("./generator/templates/to_protobuf/types/id.txt")
	if err != nil {
		panic(err)
	}

	to_protobuf_id = string(s)

	s, err = ioutil.ReadFile("./generator/templates/to_protobuf/types/key.txt")
	if err != nil {
		panic(err)
	}

	to_protobuf_key = string(s)

	s, err = ioutil.ReadFile("./generator/templates/to_protobuf/types/time_duration.txt")
	if err != nil {
		panic(err)
	}

	to_protobuf_time_duration = string(s)

}

func GenerateToProtobufs(structure Struct) string {
	s := ""

	for _, field := range structure.fields {
		if !field.config.toProtobuf {
			continue
		}

		replacer := field.Replacer()

		if replacer == nil {
			continue
		}

		ty := field.ty.String(field.config)

		switch ty {
		case "bool":
			s += replacer.Replace(to_protobuf_bool)
		case "Hbar":
			s += replacer.Replace(to_protobuf_hbar)
		case "AccountID", "FileID", "ContractID", "NftID", "TokenID", "TopicID", "ScheduleID":
			s += replacer.Replace(to_protobuf_id)
		case "Key":
			s += replacer.Replace(to_protobuf_key)
		case "time.Duration":
			s += replacer.Replace(to_protobuf_time_duration)
		case "string":
			s += replacer.Replace(to_protobuf_string)
		default:
			panic(fmt.Sprintf("Attempted to generate to protobuf conditional for type %s, but cannot find template for said type under ./generators/templates/to_protobuf/types/", ty))
		}
	}

	s = strings.ReplaceAll(to_protobuf_function, "<conditionals>", s)

	replacer := strings.NewReplacer(
		"<this.type>", structure.name,
		"<proto.type>", structure.protoName,
	)

	return replacer.Replace(s)
}
