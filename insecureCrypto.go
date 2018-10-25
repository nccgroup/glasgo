// Copyright 2018 Terence Tarvis.  All rights reserved.
//  

package main

import (
	"go/ast"
	"strings"
)

func init() {
	register("insecureCrypto",
		"this test checks for insecure cryptography primitives",
		cryptoCheck,
		fileNode)
}

func insecureCalls() map[string]bool {
	calls := make(map[string]bool)
	calls["crypto/des"] 	= true
	calls["crypto/md5"] 	= true
	calls["crypto/rc4"] 	= true
	calls["crypto/sha1"] 	= true

	return calls;
}

// getImports returns all the imports for a full file AST node
// todo: consider moving this out somewhere to a helper function to extract imports
func getImports(fn *ast.File) []string {
	var imported []string
	for _, pkg := range fn.Imports {
		// trim parantheses
		pkgName := strings.Trim(pkg.Path.Value, "\"");
		imported = append(imported, pkgName);
	}
	return imported;
}

func cryptoCheck(f *File, node ast.Node) {
	var imported []string;
	insecure := insecureCalls();

	if fileNode, ok := node.(*ast.File); ok {
		imported = getImports(fileNode);
	}
	for _, call := range imported {
		if _, ok := insecure[call]; ok {
			f.Reportf(node.Pos(), "insecure cryptographic import: %s", call);
		}
	}
	return;
}
