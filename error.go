// Copyright 2018 Terence Tarvis.  All rights reserved.

package main

import (
	"go/ast"
	"go/types"
)

func init() {
	register("error",
		"this tests to see if any errors were ignored",
		errorCheck,
		assignStmt,
		exprStmt)
}

func returnsError(f *File, call *ast.CallExpr) int {
	if typeValue := f.pkg.info.TypeOf(call); typeValue != nil {
		switch t := typeValue.(type) {
		case *types.Tuple:
			for i := 0; i < t.Len(); i++ {
				variable := t.At(i)
				if variable != nil && variable.Type().String() == "error" {
					return i;
				}
			}
		case *types.Named:
			if t.String() == "error"{
				return 0;
			}
		}	
	}
	return -1;
}
 
// Possibly check if anything returns an error before running the test
// however, this may take roughly the same amount of effort as
// just running the test in the first place.
//
func errorCheck(f *File, node ast.Node) {
	switch stmt := node.(type) {
	case *ast.AssignStmt:
		for _, rhs := range stmt.Rhs {
			if call, ok := rhs.(*ast.CallExpr); ok {
				index := returnsError(f, call)
				if index < 0 {
					continue
				}
				lhs := stmt.Lhs[index]
				if id, ok := lhs.(*ast.Ident); ok && id.Name == "_" {
					// todo real reporting
					re := f.ASTString(rhs);
					le := f.ASTString(lhs);
					f.Reportf(stmt.Pos(), "error ignored %s %s", le, re);
				}
			}
		}
	case *ast.ExprStmt:
		if expr, ok := stmt.X.(*ast.CallExpr); ok {
			pos := returnsError(f, expr);
			if pos >= 0 {
				// todo real reporting
				x := f.ASTString(expr);
				f.Reportf(stmt.Pos(), "error ignored %s", x);
			}
		}
	}
}
