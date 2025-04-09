// Analyzer checks for direct calls to os.Exit in the main function of the main package.
//
// This analyzer detects and reports direct usage of os.Exit within the main function.
package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "Analyzer that prohibits direct calls to os.Exit in the main function of the main package.",
	Run:  run,
}

// -exclude flag excluding directories for analysis.
var exclude string // -exclude flag
var excludeDoc = `excluding directories for analysis; 
you must specify the directories separated by commas; 
example: /my-dir,/other-dir`

func init() {
	ExitCheckAnalyzer.Flags.StringVar(&exclude, "exclude", exclude, excludeDoc)
}

func run(pass *analysis.Pass) (any, error) {
	excludeDirs := strings.Split(exclude, ",")

	for _, file := range pass.Files {
		for _, excludePath := range excludeDirs {
			if strings.Contains(pass.Fset.Position(file.Package).Filename, excludePath) {
				//nolint:all //declaration of the run function
				// of the field of the ExitCheckAnalyzer structure
				// requires the return of any, erorr (nilnil linter error)
				return nil, nil
			}
		}
		if file.Name.String() != "main" {
			//nolint:all //declaration of the run function
			// of the field of the ExitCheckAnalyzer structure
			// requires the return of any, erorr (nilnil linter error)
			return nil, nil
		}
		ast.Inspect(file, func(curNode ast.Node) bool {
			var mainFunc *ast.FuncDecl
			var ok bool
			if mainFunc, ok = curNode.(*ast.FuncDecl); !ok || mainFunc.Name.String() != "main" {
				return true
			}
			ast.Inspect(mainFunc, func(node ast.Node) bool {
				var expr *ast.CallExpr
				var ok bool
				if expr, ok = node.(*ast.CallExpr); !ok {
					return true
				}
				var selector *ast.SelectorExpr
				if selector, ok = expr.Fun.(*ast.SelectorExpr); !ok {
					return true
				}
				var identifier *ast.Ident
				if identifier, ok = selector.X.(*ast.Ident); !ok {
					return true
				}
				if selector.Sel.Name == "Exit" && identifier.Name == "os" {
					pass.Reportf(identifier.NamePos, "Detected direct call to os.Exit in the main function of the main package")
				}

				return true
			})

			return true
		})
	}

	//nolint:all //declaration of the run function
	// of the field of the ExitCheckAnalyzer structure
	// requires the return of any, erorr (nilnil linter error)
	return nil, nil
}
