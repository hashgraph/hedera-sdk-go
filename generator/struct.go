package main

import (
	"go/ast"
	"go/doc"
	"strings"
)

type Structs struct {
	structs []Struct
}

type Struct struct {
	name          string
	fileName      string
	fields        []Field
	protoName     string
	protoAccessor string
}

func StructsFromFiles(files map[string]*ast.File, documentation *doc.Package) []Struct {
	for fileName, file := range files {
        return StructsFromFile(fileName, file, documentation)
	}

    return make([]Struct, 0)
}

func StructsFromFile(fileName string, file *ast.File, documentation *doc.Package) []Struct {
    for _, decl := range file.Decls {
        switch decl.(type) {
        case *ast.GenDecl:
            return StructsFromGenDecl(fileName, decl.(*ast.GenDecl), documentation)
        }
    }

    return make([]Struct, 0)
}

func StructsFromDecl(fileName string, decl ast.Decl, documentation *doc.Package) []Struct {
    switch decl.(type) {
    case *ast.GenDecl:
        return StructsFromGenDecl(fileName, decl.(*ast.GenDecl), documentation)
    default:
        return make([]Struct, 0)
    }
}

func StructsFromGenDecl(fileName string, decl *ast.GenDecl, documentation *doc.Package) []Struct {
    structures := make([]Struct, 0)
    for _, spec := range decl.Specs {
        structure := StructFromSpec(fileName, spec, documentation)
        if structure == nil {
            continue
        }

        structures = append(structures, *structure)
    }

    return structures
}

func StructFromSpec(fileName string, spec ast.Spec, documentation *doc.Package) *Struct {
    switch spec.(type) {
    case *ast.TypeSpec:
        return StructFromTypeSpec(fileName, spec.(*ast.TypeSpec), documentation)
    default:
        return nil
    }
}

func StructFromTypeSpec(fileName string, spec *ast.TypeSpec, documentation *doc.Package) *Struct {
	if !((strings.HasSuffix(spec.Name.Name, "Transaction") && len(spec.Name.Name) > 12) || (strings.HasSuffix(spec.Name.Name, "Query") && len(spec.Name.Name) > 6)) {
		return nil
	}

	if spec.Name.Name != "AccountCreateTransaction" {
		return nil
	}

	d := FindDocType(documentation, spec.Name.Name)

    protoName := FindProtobufNameInDocumentation(d)
    protoAccessor := FindProtobufAccessorInDocumentation(d)
    fileName = CreateGeneratedFileName(fileName)

    if protoName == nil || protoAccessor == nil {
        return nil
    }

	return &Struct{
		name:          spec.Name.Name,
		fileName:      fileName,
		fields:        FieldsFromExpr(spec.Type),
		protoName:     *protoName,
		protoAccessor: *protoAccessor,
	}
}

func (structs Structs) String() string {
	s := ""

	for _, structure := range structs.structs {
		s += structure.String()
	}

	return s
}

func (structure Struct) String() string {
    return ""
}

// func (structure Struct) String() string {
// 	isTransaction := strings.HasSuffix(structure.name, "Transaction")
// 	param := "transaction"
// 	if !isTransaction {
// 		param = "query"
// 	}
//
// 	s := `
// package hedera
//
// import (
// 	"github.com/hashgraph/hedera-sdk-go/v2/proto"
//
// `
//
// 	if structure.usesTime {
// 		s += `
//     "time"
// `
// 	}
//
// 	s += `
// )
// `
// 	validate := ""
// 	toProtobuf := ""
// 	fromProtobuf := ""
//
// 	for _, field := range structure.fields {
// 		title := UpperInitial(field.name)
// 		protoTitle := title
// 		singular := false
// 		singularName := field.name
// 		protobufIgnored := false
// 		for _, option := range field.options {
// 			if strings.Contains(option, "=") {
// 				s := strings.Split(option, "=")
// 				switch s[0] {
// 				case "protobuf":
// 					protoTitle = UpperInitial(s[1])
// 				case "ignore":
// 					protobufIgnored = s[1] == "protobuf"
// 				}
// 			} else {
// 				switch option {
// 				case "singular":
// 					singular = true
// 					singularName = UpperInitial(field.name[0 : len(field.name)-1])
// 				}
// 			}
// 		}
//
// 		if len(field.options) > 0 && !strings.Contains(field.options[0], "=") {
// 			protoTitle = strings.Title(field.options[0])
// 		}
//
// 		if singular {
// 			s += fmt.Sprintf(`
// func (%s *%s) Set%s(%s %s) *%s {
//     %s.%s = []%s{%s}
//     return %s
// }
// `, param, structure.name, singularName, LowerInitial(singularName), field.ty[2:], structure.name, param, field.name, field.ty[2:], LowerInitial(singularName), param)
// 		} else {
// 			s += fmt.Sprintf(`
// func (%s *%s) Set%s(%s %s) *%s {
//     %s.%s = %s
//     return %s
// }
// `, param, structure.name, title, field.name, field.ty, structure.name, param, field.name, field.name, param)
// 		}
// 		switch field.ty {
// 		case "Key":
// 			s += fmt.Sprintf(`
// func (%s %s) Get%s() (%s, error) {
//     return %s.%s, nil
// }
// `, param, structure.name, title, field.ty, param, field.name)
// 		default:
// 			if singular {
// 				s += fmt.Sprintf(`
// func (%s %s) Get%s() %s {
//     if len(%s.%s) > 0 {
//         return %s.%s[0]
//     } else {
//         return %s{}
//     }
// }
// `, param, structure.name, singularName, field.ty[2:], param, field.name, param, field.name, field.ty[2:])
// 			} else {
// 				s += fmt.Sprintf(`
// func (%s %s) Get%s() %s {
//     return %s.%s
// }
// `, param, structure.name, title, field.ty, param, field.name)
// 			}
// 		}
//
// 		if protobufIgnored {
// 			continue
// 		}
//
// 		switch field.ty {
// 		case "Key":
// 			toProtobuf += fmt.Sprintf(`
//     if %s.%s != nil {
//         pb.%s = %s.%s.toProtoKey()
//     }
//     `, param, field.name, protoTitle, param, field.name)
// 		case "Hbar":
// 			toProtobuf += fmt.Sprintf(`
//     if %s.%s.tinybar != 0 {
//         pb.%s = uint64(%s.%s.tinybar)
//     }
//     `, param, field.name, protoTitle, param, field.name)
// 		case "string":
// 			toProtobuf += fmt.Sprintf(`
//     if %s.%s != "" {
//         pb.%s = %s.%s
//     }
//     `, param, field.name, protoTitle, param, field.name)
// 		case "time.Duration":
// 			toProtobuf += fmt.Sprintf(`
//     if %s.%s != 0 {
//         pb.%s = durationToProtobuf(%s.%s)
//     }
//     `, param, field.name, protoTitle, param, field.name)
// 		case "bool":
// 			toProtobuf += fmt.Sprintf(`
//     pb.%s = %s.%s
//         `, protoTitle, param, field.name)
// 		case "AccountID", "FileID", "ContractID", "TopicID", "TokenID", "ScheduleID", "NftID":
// 			toProtobuf += fmt.Sprintf(`
//     if !%s.%s.isZero() {
//         pb.%s = %s.%s.toProtobuf()
//     }
//     `, param, field.name, protoTitle, param, field.name)
// 		}
//
// 		if strings.HasSuffix(field.name, "ID") {
// 			validate += fmt.Sprintf(`
//     if !%s.%s.isZero() {
//         if err := %s.%s.Validate(client); err != nil {
//             return err
//         }
//     }
//     `, param, field.name, param, field.name)
// 		}
//
// 		switch field.ty {
// 		case "Key":
// 			fromProtobuf += fmt.Sprintf(`
//     if pb.%s != nil {
//         %s, _ := keyFromProtobuf(pb.%s, nil)
//         tx.Set%s(%s)
//     }
//             `, protoTitle, field.name, protoTitle, title, field.name)
// 		case "AccountID", "FileID", "ContractID", "TopicID", "TokenID", "ScheduleID", "NftID":
// 			fromProtobuf += fmt.Sprintf(`
//     if pb.%s != nil {
//         %s := %sFromProtobuf(pb.%s, nil)
//         tx.Set%s(%s)
//     }
//             `, protoTitle, field.name, LowerInitial(field.ty), protoTitle, title, field.name)
// 		case "Hbar":
// 			fromProtobuf += fmt.Sprintf(`
//     if pb.%s != 0 {
//         %s := HbarFromTinybar(int64(pb.%s))
//         tx.Set%s(%s)
//     }
//             `, protoTitle, field.name, protoTitle, title, field.name)
// 		case "bool":
// 			fromProtobuf += fmt.Sprintf(`
//     tx.Set%s(pb.%s)
//             `, title, protoTitle)
// 		case "time.Duration":
// 			fromProtobuf += fmt.Sprintf(`
//     if pb.%s != nil {
//         %s := durationFromProtobuf(pb.%s)
//         tx.Set%s(%s)
//     }
//             `, protoTitle, field.name, protoTitle, title, field.name)
// 		}
// 	}
//
// 	client := "client"
// 	if validate == "" {
// 		client = "_"
// 	}
//
// 	s += fmt.Sprintf(
// 		`
// func (%s %s) validateChecksums(%s *Client) error {
//     %s
//     return nil
// }
// `, param, structure.name, client,
// 		validate)
//
// 	s += fmt.Sprintf(`
// func (%s %s) toProtobuf() *proto.%s {
//     pb := &proto.%s{}
//     %s
//     return pb
// }
// `, param, structure.name, structure.protoName,
// 		structure.protoName,
// 		toProtobuf)
//
// 	if isTransaction {
// 		s += fmt.Sprintf(`
// func %sFromProtobuf(transaction Transaction, body *proto.TransactionBody) *%s {
//     tx := New%s()
//     pb := body.Get%s()
//     %s
//     return tx
// }
// `, LowerInitial(structure.name), structure.name, structure.name, structure.protoAccessor, fromProtobuf)
//
// 		s += fmt.Sprintf(`
// func (%s *%s) onFreeze(body *proto.TransactionBody, pb *proto.%s) bool {
//     body.Data = &proto.TransactionBody_%s{
//         %s: pb,
//     }
//     return true
// }
// `, param, structure.name, structure.protoName, structure.protoAccessor, structure.protoAccessor)
// 	}
//
// 	return s
// }
