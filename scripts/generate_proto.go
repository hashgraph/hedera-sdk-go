package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	// nolint:dogsled
	_, buildFilename, _, _ := runtime.Caller(0)
	projectDir := path.Join(buildFilename, "../..")

	// remove all existing files

	removeAllWithExt(projectDir, "services", ".pb.go")

	// invoke the build command for this module

	buildServices(projectDir)
}

func buildServices(dir string) {
	// collect all files from proto/services/_
	var servicesProtoFiles []string

	var servicesModuleDecls []string

	err := filepath.Walk(path.Join(dir, "services/hapi/hedera-protobufs/services"), func(filename string, info fs.FileInfo, err error) error {
		if !strings.HasSuffix(filename, ".proto") {
			return nil
		}

		pathFromRoot := strings.TrimPrefix(filename, dir+"/")
		pathBase := path.Base(filename)
		servicesProtoFiles = append(servicesProtoFiles, pathFromRoot)

		// the -M argument generation is what allows us to avoid
		// requiring "option go_package"

		servicesModuleDecls = append(servicesModuleDecls,
			fmt.Sprintf("--go_opt=M%v=github.com/hashgraph/hedera-go-sdk/services", pathBase),
			fmt.Sprintf("--go-grpc_opt=M%v=github.com/hashgraph/hedera-go-sdk/services", pathBase),
		)

		return nil
	})
	if err != nil {
		panic(err)
	}

	// generate proto files for services code

	cmdArguments := []string{
		"--go_out=proto/services/",
		"--go_opt=paths=source_relative",
		"--go-grpc_out=proto/services/",
		"--go-grpc_opt=paths=source_relative",
		"-Iservices/hapi/hedera-protobufs/services",
	}

	cmdArguments = append(cmdArguments, servicesModuleDecls...)
	cmdArguments = append(cmdArguments, servicesProtoFiles...)

	cmd := exec.Command("protoc", cmdArguments...)
	cmd.Dir = dir

	mustRunCommand(cmd)
	renamePackageDeclGrpcFiles(dir, "proto", "services")
}

// deprecated
// nolint
func buildMirror(dir string) {
	cmd := exec.Command("protoc",
		"--go_out=proto/",
		"--go_opt=Mbasic_types.proto=github.com/hashgraph/hedera-sdk-go/v2/proto/services",
		"--go_opt=Mtimestamp.proto=github.com/hashgraph/hedera-sdk-go/v2/proto/services",
		"--go_opt=Mconsensus_submit_message.proto=github.com/hashgraph/hedera-sdk-go/v2/proto/services",
		"--go_opt=Mmirror/consensus_service.proto=github.com/hashgraph/hedera-sdk-go/v2/proto/mirror",
		"--go_opt=Mmirror/mirror_network_service.proto=github.com/hashgraph/hedera-sdk-go/v2/proto/mirror",
		"--go_opt=paths=source_relative",
		"--go-grpc_out=proto/",
		"--go-grpc_opt=Mbasic_types.proto=github.com/hashgraph/hedera-sdk-go/v2/proto/services",
		"--go-grpc_opt=Mtimestamp.proto=github.com/hashgraph/hedera-sdk-go/v2/proto/services",
		"--go-grpc_opt=Mconsensus_submit_message.proto=github.com/hashgraph/hedera-sdk-go/v2/proto/services",
		"--go-grpc_opt=Mmirror/consensus_service.proto=github.com/hashgraph/hedera-sdk-go/v2/proto/mirror",
		"--go-grpc_opt=Mmirror/mirror_network_service.proto=github.com/hashgraph/hedera-sdk-go/v2/proto/mirror",
		"--go-grpc_opt=paths=source_relative",
		"--proto_path=proto/",
		"-Iproto/mirror",
		"-Iproto/services",
		"proto/mirror/consensus_service.proto",
		"proto/mirror/mirror_network_service.proto",
	)

	cmd.Dir = dir

	mustRunCommand(cmd)
	renamePackageDeclGrpcFiles(dir, "com_hedera_mirror_api_proto", "mirror")
}

// deprecated
// nolint
func buildSdk(dir string) {
	var servicesProtoFiles []string

	var servicesModuleDecls []string

	err := filepath.Walk(path.Join(dir, "services/hapi/hedera-protobufs/sdk"), func(filename string, info fs.FileInfo, err error) error {
		if !strings.HasSuffix(filename, ".proto") {
			return nil
		}

		pathFromRoot := strings.TrimPrefix(filename, dir+"/")
		pathBase := path.Base(filename)
		servicesProtoFiles = append(servicesProtoFiles, pathFromRoot)

		servicesModuleDecls = append(servicesModuleDecls,
			fmt.Sprintf("--go_opt=M%v=github.com/hashgraph/hedera-sdk-go/v2/proto/services", pathBase),
			fmt.Sprintf("--go-grpc_opt=M%v=github.com/hashgraph/hedera-sdk-go/v2/proto/services", pathBase),
		)

		return nil
	})
	if err != nil {
		panic(err)
	}

	cmdArguments := []string{
		"--go_out=proto/sdk/",
		"--go_opt=Mtransaction_list.proto=github.com/hashgraph/hedera-sdk-go/v2/proto/sdk",
		"--go_opt=paths=source_relative",
		"--go-grpc_out=proto/sdk/",
		"--go-grpc_opt=Mtransaction_list.proto=github.com/hashgraph/hedera-sdk-go/v2/proto/sdk",
		"--go-grpc_opt=paths=source_relative",
		"-Iproto/sdk",
		"-Iproto/services",
	}

	cmdArguments = append(cmdArguments, servicesModuleDecls...)
	cmdArguments = append(cmdArguments, "proto/sdk/transaction_list.proto")

	cmd := exec.Command("protoc", cmdArguments...)
	cmd.Dir = dir

	mustRunCommand(cmd)
}

func mustRunCommand(cmd *exec.Cmd) {
	_, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			fmt.Print(string(exitErr.Stderr))
			os.Exit(exitErr.ExitCode())
		}

		panic(err)
	}
}

func removeAllWithExt(dir string, module string, ext string) {
	err := filepath.Walk(path.Join(dir, module), func(filename string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(filename, ext) {
			err := os.Remove(filename)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func renamePackageDeclGrpcFiles(dir string, oldPackage string, newPackage string) {
	err := filepath.Walk(path.Join(dir, newPackage), func(filename string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(filename, "_grpc.pb.go") {
			data, err := os.ReadFile(filename)
			if err != nil {
				return err
			}

			contents := string(data)
			contents = strings.Replace(contents,
				fmt.Sprintf("package %s", oldPackage),
				fmt.Sprintf("package %s", newPackage),
				1,
			)

			return os.WriteFile(filename, []byte(contents), info.Mode())
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}
