// Copyright 2018 Terence Tarvis.  All rights reserved.
// add a license

package main

import (
	"go/ast"
)

func init() {
	register("closeCheck",
		"this tests if things with .Close() method have .Close() actually called on them",
		closeCheck,
		funcDecl)
}

func opensFile(f *File, x ast.Expr) bool {
	if(f.pkg.info.TypeOf(x).String() == "(*os.File, error)") {
		return true
	}
	return false;
}

// closesFile checks the remaining statements in a function body for a .Close() method
func closesFile(f *File, stmts []ast.Stmt) bool {
	for _, stmt := range stmts {
		switch expr := stmt.(type) {
		case *ast.AssignStmt:
			rhs := expr.Rhs;
			for _, x := range rhs {
				name, err := getFullFuncName(x);
				if err != nil {
					warnf("issue, %v", err);
				}
				if(name == "file/Close") {
					return true
				}
			}
		case *ast.ExprStmt:
			name, err := getFullFuncName(expr.X);
			if err != nil {
				warnf("issue, %v", err);
			}
			if(name == "file/Close") {
				return true
			}
		}
	}
	return false
}

// for the time being this just checks a function to see if an opened file is closed
// http.MaxBytesReader should also be checked for a close
func closeCheck(f *File, node ast.Node) {
	var formatString string = "Audit for Close() method called on opened file, %s"
	// loop through block
	// look for file open
	// look for file close
	// if no file close, report
	// consider checking if the opened file is returned
	// consider checking if an open file is an input and no file is returned
	// it turns out you can open a file and not use it. What then?
	// ugh I really hate this
	// is walking the statements a better option?
	if fun, ok := node.(*ast.FuncDecl); ok {
		for i, stmts := range fun.Body.List {
			switch stmt := stmts.(type) {
			case *ast.AssignStmt:
				rhs := stmt.Rhs;
				for _, x := range rhs {
					if(opensFile(f, x)) {
						if(!closesFile(f, fun.Body.List[i:])) {
							f.Reportf(stmt.Pos(), formatString,f.ASTString(x))
						}
					}
				}
			case *ast.ExprStmt:
				if(opensFile(f, stmt.X)) {
					if(!closesFile(f, fun.Body.List[i:])) {
						f.Reportf(stmt.Pos(), formatString, f.ASTString(stmt.X))
					}
				}
			case *ast.IfStmt:
				if s, ok := stmt.Init.(*ast.AssignStmt); ok {
					rhs := s.Rhs;
					for _, x := range rhs {
						if(opensFile(f, x )) {
							if(!closesFile(f, fun.Body.List[i:])) {
								f.Reportf(stmt.Pos(), formatString, f.ASTString(x))
							}
						}
					}
				}
			default:
				// do nothing for time being
			}
		}
	}
	return;
}
