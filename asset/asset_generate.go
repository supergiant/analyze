// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/httpfs/union"
	"github.com/shurcooL/vfsgen"
)

func main() {

	var fileSystem http.FileSystem = union.New(map[string]http.FileSystem{
		"/ui":      http.Dir("./ui"),
		"/swagger": http.Dir("./swagger"),
	})

	err := vfsgen.Generate(fileSystem, vfsgen.Options{
		PackageName:  "asset",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
