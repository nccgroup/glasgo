// Copyright 2018 Terence Tarvis.  All rights reserved.
// add a license

package main

import (
	"go/ast"
)

func init() {
	register("textTemp",
		"this is a test to see if template/text and http methods are in use",
		textTempCheck,
		fileNode)
}

func textTempCheck(f *File, node ast.Node) {
	importedPkgs := make(map[string]bool);
	if fileNode, ok := node.(*ast.File); ok {
		imports := getImports(fileNode);
		for _, imported := range imports {
			importedPkgs[imported] = true;	
		}		
		if a, b := importedPkgs["net/http"], importedPkgs["text/template"]; a && b {
			f.Reportf(node.Pos(), "audit use of text/template in HTTP responses");
		}
	}

	return;
}
