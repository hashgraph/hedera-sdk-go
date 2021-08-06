package main

import (
	"fmt"
	"strings"
)

var to_protobuf_function string = ReadFileToString("./generator/templates/to_protobuf/function.txt")
var to_protobuf_bool string = ReadFileToString("./generator/templates/to_protobuf/types/bool.txt")
var to_protobuf_hbar string = ReadFileToString("./generator/templates/to_protobuf/types/hbar.txt")
var to_protobuf_id string = ReadFileToString("./generator/templates/to_protobuf/types/id.txt")
var to_protobuf_key string = ReadFileToString("./generator/templates/to_protobuf/types/key.txt")
var to_protobuf_time_duration string = ReadFileToString("./generator/templates/to_protobuf/types/time_duration.txt")
var to_protobuf_string string = ReadFileToString("./generator/templates/to_protobuf/types/string.txt")

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
