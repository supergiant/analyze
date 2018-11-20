// asset package provides the assets via a virtual filesystem.
package asset

//package main

import (
	// The blank import is to make govendor happy.
	_ "github.com/shurcooL/vfsgen"
)

//go:generate go run asset_generate.go
