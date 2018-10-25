// Copyright 2018 Terence Tarvis.  All rights reserved.
//  

package main

import (
	"go/ast"
)

func init() {
	register("insecureRand",
		"this is test to check if random nums generated insecurely",
		randCheck,
		fileNode)
}

func randCheck(f *File, node ast.Node) {
	var imported []string
	
	if fileNode, ok := node.(*ast.File); ok {
		imported = getImports(fileNode);
	}

	for _, pkg := range imported {
		if(pkg == "math/rand") {
			f.Reportf(node.Pos(), "audit the use of insecure random number generator: import: %s", pkg);
		} 
	}
	return;
}
