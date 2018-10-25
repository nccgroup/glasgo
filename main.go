package main

import (
	"fmt"
	"flag"
	"go/ast"
	"go/build"
	"go/token"
	"go/parser"
	"go/printer"
	"go/types"
	"go/importer"
	"bytes"
	"strings"
	"os"
	"path/filepath"
)

var stdImporter types.Importer

var (
	source = flag.Bool("source", false, "import from source instead of compiled object files")
)

// a global variable for the exit code.
var exitCode = 0;

var report = make(map[string]bool);

var (
	// shortens type names
	// These are the relevant AST node types to check
	// with corresponding cases
	assignStmt	*ast.AssignStmt
	binaryExpr	*ast.BinaryExpr
	callExpr	*ast.CallExpr
	compositeLit	*ast.CompositeLit
	exprStmt	*ast.ExprStmt
	fileNode	*ast.File
	forStmt		*ast.ForStmt
	funcDecl	*ast.FuncDecl
	funcLit		*ast.FuncLit
	genDecl		*ast.GenDecl
	interfaceType	*ast.InterfaceType
	rangeStmt	*ast.RangeStmt
	returnStmt	*ast.ReturnStmt
	structType	*ast.StructType
)

var (
	// checkers is a map to a map
	// the map maps AST types to maps of checker names to checker functions
	// this is to first get the functions needed for a certain type
	// and second to take just the functions we want to run.
	// refactor this to map to struct
	checkers	= make(map[ast.Node]map[string]func(*File, ast.Node))
	
)

// A map 
// File is a visitor type for the parse tree.
// it also contains the corresponding AST to a parsed file
// pkg contains data on the entire package that was parsed
// this includes things like type info so you can spot
// an expression, like a func call, and look up it's type
type File struct {
	pkg	*Package
	fset	*token.FileSet
	name	string
	file	*ast.File

	b	bytes.Buffer // used for logging and printing results

	// a map of all registered checkers to run for each node
	checkers map[ast.Node][]func(*File, ast.Node);
}

// Reportf reports issues to a log for each file for later printing
func (f *File) Reportf(pos token.Pos, format string, args ...interface{}) {
	// update this to use a logger
	fmt.Fprintf(os.Stderr, "\t* %v %s \n", f.loc(pos), fmt.Sprintf(format, args...));
}

// loc (line of code) returns a formatted string of file and a file position
func (f *File) loc(pos token.Pos) string {
	if pos == token.NoPos {
		return ""
	}
	// we won't print column, just line
	posn := f.fset.Position(pos)
	return fmt.Sprintf("%s:%d", posn.Filename, posn.Line);
}

// warnf is a formatted error printer that does not exit
// but it does set an exit code.
func warnf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "{insert tool name here}: "+format+"\n", args...);
	exitCode = 1;
}

// register registers the named checker function
// to be called with AST nodes of the given types.
func register(name, usage string, fn func(*File, ast.Node), types ...ast.Node) {
	report[name] = true;
	for _, typ := range types {
		m, ok := checkers[typ];
		if !ok {
			m = make(map[string]func(*File, ast.Node));
			checkers[typ] = m;
		}
		m[name] = fn;
	}
}

// Visit implements the visitor interface we need to walk the tree
// ast.Walk calls v.Visit(node)
func (f *File) Visit(node ast.Node) ast.Visitor {
	var key ast.Node
	switch node.(type) {
	case *ast.AssignStmt:
		key = assignStmt
	case *ast.BinaryExpr:
		key = binaryExpr
	case *ast.CallExpr:
		key = callExpr
	case *ast.CompositeLit:
		key = compositeLit
	case *ast.ExprStmt:
		key = exprStmt
	case *ast.File:
		key = fileNode
	case *ast.ForStmt:
		key = forStmt
	case *ast.FuncDecl:
		key = funcDecl
	case *ast.FuncLit:
		key = funcLit
	case *ast.GenDecl:
		key = genDecl
	case *ast.InterfaceType:
		key = interfaceType
	case *ast.RangeStmt:
		key = rangeStmt
	case *ast.ReturnStmt:
		key = returnStmt
	case *ast.StructType:
		key = structType
	}
	// runs checkers below
	for _, fn := range f.checkers[key] {
		fn(f, node)
	}
	return f;
}

type Package struct {
	path	string
	types 	map[ast.Expr]types.TypeAndValue;
	typePkg	*types.Package
	info	*types.Info
}

func (pkg *Package) check(fs *token.FileSet, astFiles []*ast.File) error {
	if stdImporter == nil {
		if *source {
			stdImporter = importer.For("source", nil)
		} else {
			stdImporter = importer.Default();
		}
	}
	pkg.types = make(map[ast.Expr]types.TypeAndValue);

	conf := types.Config{
		Importer: stdImporter,
		Error: func(err error) { 
				// todo refactor this
				fmt.Printf("\tWarning: during type checking, %v\n", err)
			},
	}

	info := types.Info{
		Types: pkg.types,
	}

	// Type-Check the package.
	typePkg, err := conf.Check(pkg.path, fs, astFiles, &info);
	pkg.typePkg = typePkg
	pkg.info = &info;
	return err;
	
}

// checkPackageDir extracts the go files from a directory and passes them to 
// checkPackage for analysis
func checkPackageDir(directory string) {
	context := build.Default
	// gets build tags if any exist in order to preserve them through the coming import
	/*
	these are commented out until proof is made of being necessary
	if len(context.BuildTags) != 0 {
		warnf("build tags already set: %s," context.BuildTags);
	}
	context.BuildTags = append(tagList, context.BuildTags...);
	*/

	pkg, err := context.ImportDir(directory, 0); // 0 means no ImportMode is set i.e. default
	if err != nil {
		// no go source files
		if _, noGoSource := err.(*build.NoGoError); noGoSource {
			return;
		}
		// not considered fatal because we are recursively walking directories
		warnf("error processing directory %s, %s", directory, err);
		return;
	}
	var names []string
	names = append(names, pkg.GoFiles...);
	names = append(names, pkg.CgoFiles...);
	names = append(names, pkg.TestGoFiles...);
	/* there are other types include binary files that can be added */
	
	/* prefix each file with the directory name
	 * could use a refactor
	*/
	if directory != "." {
		for i, name := range names{
			names[i] = filepath.Join(directory, name);
		}
	}
	checkPackage(names);
}

// checkPackage runs analysis on all named files in a package.
// It parses and then runs the analysis.
// It returns the parsed package or nil.
func checkPackage(names []string) {
	var files []*File;
	var astFiles []*ast.File;
	fset := token.NewFileSet();
	var err error;
	for _, name := range names {
		// skipping using ioutil to read the file data
		// and just going to parse files directly.
		var parsedFile *ast.File;
		if strings.HasSuffix(name, ".go") {
			parsedFile, err = parser.ParseFile(fset, name, nil, parser.ParseComments)
			if err != nil {
				// warn but continue
				warnf("error: %s: %s", name, err);
				return;
			}
			astFiles = append(astFiles, parsedFile);
		}
		file := &File{
			fset:	fset,
			name:	name,
			file:	parsedFile,
		}
		files = append(files, file);
	}
	if len(astFiles) == 0 {
		return;
	}
	pkg := new(Package);
	
	// Type check package and
	// generate information about it
	err = pkg.check(fset, astFiles);
	if err != nil {
		// probably should just keep going
		// fmt.Printf("exited, %v", err);
		//os.Exit(0);
		// errors being caught in different location.
	}

	// Check.
	for _, file := range files {
		file.pkg = pkg;
	}

	chk := make(map[ast.Node][]func(*File, ast.Node));
	for typ, set := range checkers {
		for name, fn := range set {
			// check to see if named function will be run and reported
			_, ok := report[name];
			if ok {
				chk[typ] = append(chk[typ], fn);
			}
		}
	}
	for _, file := range files {
		file.checkers = chk
		if file.file != nil {
			// Should this go in to a new function to make it more readable?
			// file.walkFile(file.name, file.file) as a method?
			fmt.Printf("Checking %s\n", file.name);
			ast.Walk(file, file.file);
		}
	}
}

// visit is for walking input directory roots
func visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		warnf("directory walk error: %s", err);
		return err;
	}
	// make sure we are only dealing with directories here
	if !info.IsDir() {
		return nil
	}
	checkPackageDir(path);
	return nil;
}

// ASTString returns a string representation of the AST for reporting
func (f *File) ASTString(x ast.Expr) string {
	var b bytes.Buffer
	printer.Fprint(&b, f.fset, x);
	return b.String()
}

// getFuncName returns just function name i.e. not ioutil.ReadAll but just ReadAll
// not returning errors,
func getFuncName(node ast.Node) string {
	if call, ok := node.(*ast.CallExpr); ok {
		if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
			if(fun.Sel.Name != "") {
				return fun.Sel.Name;
			}
		}
		if fun, ok := call.Fun.(*ast.Ident); ok {
			if(fun.Name != "") {
				return fun.Name;
			}
		}
	} 
	return ""
}

// getFullFuncName extracts a full function name path i.e ioutil.ReadAll
func getFullFuncName(node ast.Node) (string, error) {
	var names []string
	var callName string
	if call, ok := node.(*ast.CallExpr); ok {
		if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
			// fmt.Println(fun.X);
			// I think the above can be removed
			// SelectorExpr has two fields
			// X and Sel
                        // X (through reflection) was found to be an Ident
                        // Sel has field Name
                        // Ident's have a field Name also.
			if id, ok := (fun.X).(*ast.Ident); ok {
				names = append(names, id.Name);
				names = append(names, fun.Sel.Name);
				callName = strings.Join(names, "/")
				return callName, nil
			}
		}
	}
	return "", fmt.Errorf("type conversion of CallExpr failed, no name extracted, %v", node);
}

func main() {
	var runOnDirs, runOnFiles bool;
	flag.Parse();

	for _, name := range flag.Args() {
		// check to see if cl argument is a directory
		f, err := os.Stat(name);
		if err != nil {
			warnf("error: %s", err);
			continue;
		}
		if f.IsDir() {
			runOnDirs = true;
		} else {
			runOnFiles = true;
		}
	}
	if runOnDirs && runOnFiles {
		// print an error
		fmt.Println("error: input arguments must not be both directories and files");
		exitCode = 1;
		os.Exit(exitCode);
	}
	if runOnDirs {
		// I want to do each directory in order
		// so I am going to loop through these regardless
		// root is a name of a directory, at the root, to be walked
		for _, root := range flag.Args() {
			filepath.Walk(root, visit);
		}
		os.Exit(exitCode);
	}
	// else they are just file names
	fileNames := flag.Args();	
	checkPackage(fileNames);
	return;
}

