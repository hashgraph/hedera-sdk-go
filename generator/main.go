package main

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/format"
	"go/token"
	"io/ioutil"
	"path/filepath"
	// "go/types"
	"go/parser"
	// "golang.org/x/tools/go/packages"
	"regexp"
	"strings"
	"unicode"
)

var protobufTypeRegex = regexp.MustCompile("ProtobufType: .+")
var protobufAccessorRegex = regexp.MustCompile("ProtobufAccessor: .+")

type Structs struct {
	structs []Struct
}

type Struct struct {
	name      string
	fileName  string
	fields    []Field
	protoName string
	protoAccessor string
	usesTime  bool
}

type Field struct {
	name    string
	ty      string
	options []string
}

func main() {
	dir := filepath.Dir(".")
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, "./", nil, 0)
	if err != nil {
		panic(err)
	}

	// Used for AST
	pkg := pkgs["hedera"]

	pkgs, err = parser.ParseDir(fset, "./", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// Used for Comments
	// The reason we need this is because `doc.New()` takes over the `pkg` and completely removes the AST
	docPkg := pkgs["hedera"]

	structures := Structs{
		structs: make([]Struct, 0),
	}

	documentation := doc.New(docPkg, "./", 0)
	for fileName, file := range pkg.Files {
		for _, decl := range file.Decls {
			switch decl.(type) {
			case *ast.GenDecl:
				for _, spec := range decl.(*ast.GenDecl).Specs {
					switch spec.(type) {
					case *ast.TypeSpec:
						typeSpec := spec.(*ast.TypeSpec)

						if !((strings.HasSuffix(typeSpec.Name.Name, "Transaction") && len(typeSpec.Name.Name) > 12) || (strings.HasSuffix(typeSpec.Name.Name, "Query") && len(typeSpec.Name.Name) > 6)) {
							continue
						}

						if typeSpec.Name.Name != "AccountCreateTransaction" {
							continue
						}

						d := FindDocType(documentation, typeSpec.Name.Name)
						if d == "" {
							panic(fmt.Sprintf("Unable to find documentation for type name %s", typeSpec.Name.Name))
						}

						matches := protobufTypeRegex.FindAllString(d, -1)
						if matches == nil {
							panic(fmt.Sprintf(`Unable to find "ProtobufType: .*" in documentation for type name %s`, typeSpec.Name.Name))
						}

						protoName := strings.Split(matches[0], " ")[1]

						matches = protobufAccessorRegex.FindAllString(d, -1)
						if matches == nil {
							panic(fmt.Sprintf(`Unable to find "ProtobufAcessor: .*" in documentation for type name %s`, typeSpec.Name.Name))
						}

						protoAccessor := strings.Split(matches[0], " ")[1]

						fileName = strings.Split(fileName, ".go")[0]
						fileName = fmt.Sprintf("%s_accessors.go", fileName)

						structure := Struct{
							name:      typeSpec.Name.Name,
							fileName:  fileName,
							fields:    make([]Field, 0),
							protoName: protoName,
                            protoAccessor: protoAccessor,
						}

						switch typeSpec.Type.(type) {
						case *ast.StructType:
							structType := typeSpec.Type.(*ast.StructType)

							for _, field := range structType.Fields.List {
								options := make([]string, 0)
								if field.Tag != nil && field.Tag.Kind == token.STRING {
									tags := strings.Split(field.Tag.Value, " ")
									for _, t := range tags {
										if strings.HasPrefix(t, "`hedera:") {
											options = strings.Split(strings.Split(t, "\"")[1], ",")
											break
										}
									}
								}

								ty := ""

								switch field.Type.(type) {
								case *ast.Ident:
									ty = field.Type.(*ast.Ident).Name
								case *ast.StarExpr:
									starExpr := field.Type.(*ast.StarExpr)
									switch starExpr.X.(type) {
									case *ast.SelectorExpr:
										selectorExpr := starExpr.X.(*ast.SelectorExpr)
										switch selectorExpr.X.(type) {
										case *ast.Ident:
											ty = fmt.Sprintf("*%s.%s", selectorExpr.X, selectorExpr.Sel)
										default:
											panic(fmt.Sprintf("Unhandled field type: %T", selectorExpr.X))
										}
									default:
										panic(fmt.Sprintf("Unhandled field type: %T", starExpr.X))
									}
								case *ast.SelectorExpr:
									selectorExpr := field.Type.(*ast.SelectorExpr)
									switch field.Type.(*ast.SelectorExpr).X.(type) {
									case *ast.Ident:
										ty = fmt.Sprintf("%s.%s", selectorExpr.X, selectorExpr.Sel)
									default:
										panic(fmt.Sprintf("Unhandled field type: %T", field.Type.(*ast.SelectorExpr).X))
									}
								default:
									panic(fmt.Sprintf("Unhandled field type: %T", field.Type))
								}

								for _, name := range field.Names {
									if name.Name != "pb" {
										structure.fields = append(structure.fields, Field{
											name:    name.Name,
											ty:      ty,
											options: options,
										})
									}
								}

								structure.usesTime = structure.usesTime || strings.Contains(ty, "time.")
							}
						}

						structures.structs = append(structures.structs, structure)
					}
				}
			}
		}
	}

	for _, structure := range structures.structs {
		data := []byte(structure.String())
		src, err := format.Source(data)
		if err != nil {
			panic(err)
		}

		output := filepath.Join(dir, structure.fileName)
		err = ioutil.WriteFile(output, src, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func (structure Struct) String() string {
	isTransaction := strings.HasSuffix(structure.name, "Transaction")
	param := "transaction"
	if !isTransaction {
		param = "query"
	}

	s := `
package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

`

	if structure.usesTime {
		s += `
    "time"
`
	}

	s += `
)
`
	validate := ""
	toProtobuf := ""
	fromProtobuf := ""

	for _, field := range structure.fields {
		title := UpperInitial(field.name)
		protoTitle := title
		if len(field.options) > 0 && !strings.Contains(field.options[0], "=") {
			protoTitle = strings.Title(field.options[0])
		}

		s += fmt.Sprintf(`
func (%s *%s) Set%s(%s %s) *%s {
    %s.%s = %s
    return %s
}
`, param, structure.name, title, field.name, field.ty, structure.name, param, field.name, field.name, param)
		switch field.ty {
		case "Key":
			s += fmt.Sprintf(`
func (%s %s) Get%s() (%s, error) {
    return %s.%s, nil
}
`, param, structure.name, title, field.ty, param, field.name)
		default:
			s += fmt.Sprintf(`
func (%s %s) Get%s() %s {
    return %s.%s
}
`, param, structure.name, title, field.ty, param, field.name)
		}

		switch field.ty {
		case "Key":
			toProtobuf += fmt.Sprintf(`
    if %s.%s != nil {
        pb.%s = %s.%s.toProtoKey()
    }
    `, param, field.name, protoTitle, param, field.name)
		case "Hbar":
			toProtobuf += fmt.Sprintf(`
    if %s.%s.tinybar != 0 {
        pb.%s = uint64(%s.%s.tinybar)
    }
    `, param, field.name, protoTitle, param, field.name)
		case "string":
			toProtobuf += fmt.Sprintf(`
    if %s.%s != "" {
        pb.%s = %s.%s
    }
    `, param, field.name, protoTitle, param, field.name)
		case "time.Duration":
			toProtobuf += fmt.Sprintf(`
    if %s.%s != 0 {
        pb.%s = durationToProtobuf(%s.%s)
    }
    `, param, field.name, protoTitle, param, field.name)
		case "bool":
			toProtobuf += fmt.Sprintf(`
        pb.%s = %s.%s
        `, protoTitle, param, field.name)
		case "AccountID", "FileID", "ContractID", "TopicID", "TokenID", "ScheduleID", "NftID":
			toProtobuf += fmt.Sprintf(`
    if !%s.%s.isZero() {
        pb.%s = %s.%s.toProtobuf()
    }
    `, param, field.name, protoTitle, param, field.name)
		}

		if strings.HasSuffix(field.name, "ID") {
			validate += fmt.Sprintf(`
    if !%s.%s.isZero() {
        if err := %s.%s.Validate(client); err != nil {
            return err
        }
    }
    `, param, field.name, param, field.name)
		}

		switch field.ty {
		case "Key":
			fromProtobuf += fmt.Sprintf(`
    if pb.%s != nil {
        %s, _ := keyFromProtobuf(pb.%s, nil)
        tx.Set%s(%s)
    }
            `, protoTitle, field.name, protoTitle, title, field.name)
		case "AccountID", "FileID", "ContractID", "TopicID", "TokenID", "ScheduleID", "NftID":
			fromProtobuf += fmt.Sprintf(`
    if pb.%s != nil {
        %s := %sFromProtobuf(pb.%s, nil)
        tx.Set%s(%s)
    }
            `, protoTitle, field.name, LowerInitial(field.ty), protoTitle, title, field.name)
		case "Hbar":
			fromProtobuf += fmt.Sprintf(`
    if pb.%s != 0 {
        %s := HbarFromTinybar(int64(pb.%s))
        tx.Set%s(%s)
    }
            `, protoTitle, field.name, protoTitle, title, field.name)
		case "bool":
			fromProtobuf += fmt.Sprintf(`
    tx.Set%s(pb.%s)
            `, title, protoTitle)
		case "time.Duration":
			fromProtobuf += fmt.Sprintf(`
    if pb.%s != nil {
        %s := durationFromProtobuf(pb.%s)
        tx.Set%s(%s)
    }
            `, protoTitle, field.name, protoTitle, title, field.name)
		}
	}

	client := "client"
	if validate == "" {
		client = "_"
	}

	s += fmt.Sprintf(
		`
func (%s %s) validateChecksums(%s *Client) error {
    %s
    return nil
}
`, param, structure.name, client,
		validate)

	s += fmt.Sprintf(`
func (%s %s) toProtobuf() *proto.%s {
    pb := &proto.%s{}
    %s
    return pb
}
`, param, structure.name, structure.protoName,
		structure.protoName,
		toProtobuf)

	if isTransaction {
		s += fmt.Sprintf(`
func %sFromProtobuf(transaction Transaction, body *proto.TransactionBody) *%s {
    tx := New%s()
    pb := body.Get%s()
    %s
    return tx
}
`, LowerInitial(structure.name), structure.name, structure.name, structure.protoAccessor, fromProtobuf)
	}
	return s
}

func (structs Structs) String() string {
	s := ""

	for _, structure := range structs.structs {
		s += structure.String()
	}

	return s
}

func FindDocType(documentation *doc.Package, typeName string) string {
	for _, typ := range documentation.Types {
		if typ.Name == typeName {
			return typ.Doc
		}
	}
	return ""
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
