// Copyright 2018 Terence Tarvis.  All rights reserved.
// add a license

package main

import (
	"go/ast"
	"strings"
)

func init() {
	register("readAll",
		"this tests checks of use of ioutil.ReadAll needs to be audited",
		readAllCheck,
		callExpr)
}

// this checks for the bad function
// but maybe abstract this and create a function name extractor
// then throw all the bad functions together or check for them
// using it
func readAllCheck(f *File, node ast.Node) {
	// the names of the called functions
	// new function getFullFuncName does all this work.
	// todo: replace
	var names []string
	var callName string	
	if call, ok := node.(*ast.CallExpr); ok {
		if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
			// SelectorExpr has two fields
			// X and Sel
			// X (through reflection) was found to be an Ident
			// Sel has field Name
			// Ident's have a field Name also.
			if id, ok := (fun.X).(*ast.Ident); ok {
				names = append(names, id.Name);
				names = append(names, fun.Sel.Name);	
				callName = strings.Join(names, "/")
				if(callName == "ioutil/ReadAll") {
					callStr := f.ASTString(call);
					f.Reportf(node.Pos(), "audit use of ioutil.ReadAll %s", callStr);
				} 	
			}
		}
	}
	return;
}
