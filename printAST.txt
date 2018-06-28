package main

import(
	"go/ast"
	"go/token"
	"go/parser"
	"fmt"
	"reflect"
)

type visitor int;

func (v visitor) Visit(n ast.Node) ast.Visitor {
	fmt.Println(reflect.TypeOf(n));	
	return v;
}

func main() {
	var filename string = "printAST.go";
	var v visitor;

	fset := token.NewFileSet();
	file, err := parser.ParseFile(fset, filename, nil, 0);
	if err != nil {
		fmt.Printf("error, in main, %v", err);
	}

	ast.Walk(v, file); 

}
