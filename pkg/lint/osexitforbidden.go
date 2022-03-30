package lint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const (
	importSpecPathOS      = `"os"`
	defaultOsPackageAlias = "os"
	packageMain           = "main"
	functionMain          = "main"
	functionExit          = "Exit"
)

var OsExitForbiddenAnalyzer = &analysis.Analyzer{
	Name: "osexitforbidden",
	Doc:  "os.Exit call is forbidden in main func of main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	var osImportAlias string

	for _, file := range pass.Files {
		if file.Name.Name != packageMain {
			continue
		}

		osImportAlias = defaultOsPackageAlias
		ast.Inspect(file, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.ImportSpec:
				if n.Path.Value == importSpecPathOS && n.Name != nil {
					osImportAlias = n.Name.Name
				}

			case *ast.FuncDecl:
				if n.Name.Name == functionMain {
					for _, stmt := range n.Body.List {
						if isCallStmt(stmt, osImportAlias, functionExit) {
							pass.Report(analysis.Diagnostic{
								Pos:     stmt.Pos(),
								End:     stmt.End(),
								Message: "os.Exit call is forbidden in main func of main package",
							})
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func isCallStmt(stmt ast.Stmt, alias string, name string) bool {
	if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
		if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := selectorExpr.X.(*ast.Ident); ok {
					return ident.Name == alias && selectorExpr.Sel.Name == name
				}
			}
		}
	}
	return false
}
