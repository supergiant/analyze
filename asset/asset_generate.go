// +build ignore

package main

import (
	"github.com/shurcooL/vfsgen"
	"log"
	"net/http"
)

func main() {
	var fs http.FileSystem = http.Dir("./swagger")

	err := vfsgen.Generate(fs, vfsgen.Options{
		PackageName:     "asset",
		VariableName:    "SwaggerAssets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
