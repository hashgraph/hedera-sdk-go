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
