// Copyright 2018 Terence Tarvis.  All rights reserved.
// add a license

package main

import (
	"fmt"
	"go/ast"
)

func init() {
	register("squib",
		"this is just a test that does nothing",
		squibCheck,
		funcDecl, interfaceType)
}

func squibCheck(f *File, node ast.Node) {
	fmt.Println("squib");
}
