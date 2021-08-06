package main

import (
	"fmt"
	"strings"
)

var from_protobuf_function string = ReadFileToString("./generator/templates/from_protobuf/function.txt")
var from_protobuf_bool string = ReadFileToString("./generator/templates/from_protobuf/types/bool.txt")
var from_protobuf_hbar string = ReadFileToString("./generator/templates/from_protobuf/types/hbar.txt")
var from_protobuf_id string = ReadFileToString("./generator/templates/from_protobuf/types/id.txt")
var from_protobuf_key string = ReadFileToString("./generator/templates/from_protobuf/types/key.txt")
var from_protobuf_time_duration string = ReadFileToString("./generator/templates/from_protobuf/types/time_duration.txt")
var from_protobuf_string string = ReadFileToString("./generator/templates/from_protobuf/types/string.txt")

func GenerateFromProtobufs(structure Struct) string {
	s := ""

	for _, field := range structure.fields {
		if !field.config.fromProtobuf {
			continue
		}

		replacer := field.Replacer()

		if replacer == nil {
			continue
		}

		ty := field.ty.String(field.config)

		switch ty {
		case "bool":
			s += replacer.Replace(from_protobuf_bool)
		case "Hbar":
			s += replacer.Replace(from_protobuf_hbar)
		case "AccountID", "FileID", "ContractID", "NftID", "TokenID", "TopicID", "ScheduleID":
			s += replacer.Replace(from_protobuf_id)
		case "Key":
			s += replacer.Replace(from_protobuf_key)
		case "time.Duration":
			s += replacer.Replace(from_protobuf_time_duration)
		case "string":
			s += replacer.Replace(from_protobuf_string)
		default:
			panic(fmt.Sprintf("Attempted to generate from protobuf conditional for type %s, but cannot find template for said type under ./generators/templates/from_protobuf/types/", ty))
		}
	}

	s = strings.ReplaceAll(from_protobuf_function, "<conditinals>", s)

	replacer := strings.NewReplacer(
		"<this.type>", structure.name,
		"<proto.type>", structure.protoName,
	)

	return replacer.Replace(s)
}
