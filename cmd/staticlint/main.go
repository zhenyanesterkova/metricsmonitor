// This package is a multichecker - tool for static code analysis in Go.
//
// Multichecker provides a set of checks to ensure code quality.
// It can be used as part of the development workflow to identify
// potential issues in the source code.
//
// Running the tool:
//
// There are two ways to run Multichecker:
//
// 1. Via Makefile:
//   - Run the command `make staticlint-build` to build
//   - Then run `make staticlint` to start the checks
//
// 2. Directly:
//   - Build the package: `go build ./cmd/staticlint`
//   - Run with necessary parameters:
//
// `./staticlint.exe -exitcheck.exclude="/dir-name/" ./...`
//
// Running parameters:
// `-exitcheck.exclude="/dir-name/, dir-name1, /dir-name2/dir-name3"` - allows excluding certain checks
// from the analysis process.
//
// Scope of application:
// The tool recursively analyzes all packages in the current directory
// (specified via `./...`).
//
// It includes:
//
// - Standard analyzers from golang.org/x/tools/go/analysis/passes
//
//   - astalias: detects type alias usage in AST
//
//   - assign: checks for incorrect assignments
//
//   - atomic: finds issues with atomic operations
//
//   - bidirimports: checks for bidirectional imports
//
//   - buildtags: analyzes build tags
//
//   - composite: detects unnecessary composite literals
//
//   - ctxchecks: checks context usage
//
//   - deadcode: finds unused code
//
//   - depgraph: analyzes package dependencies
//
//   - errorsas: checks errors.As usage
//
//   - exportlocal: detects exported locals
//
//   - fieldalignment: optimizes struct field alignment
//
//   - govet: basic code checks
//
//   - httpresponse: checks http.ResponseWriter usage
//
//   - interfacetypes: analyzes interface type assertions
//
//   - nilness: tracks potential nil values
//
//   - printf: checks printf-style format strings
//
//   - shadow: detects variable shadowing
//
//   - simplespread: checks slice spreading
//
//   - structtag: validates struct tags
//
//   - typeparam: checks generics
//
//   - unused: finds unused identifiers
//
//   - varcheck: checks unused variables
//
// - Analyzers from staticcheck.io.
// Provides an extended set of checks including:
//
//   - Type and assertion checks
//
//   - Performance analysis
//
//   - Potential bug detection
//
//   - Style checks
//
//   - Package usage optimization
//
// - Errcheck analyzer from github.com/kisielk/errcheck/errcheck.
// Checks error handling in the code:
//
//   - Detects unchecked errors
//
//   - Analyzes function results
//
//   - Verifies error handling correctness
//
// - exitcheck analyzer detects and reports direct usage of os.Exit within the main function.
package main

import (
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"golang.org/x/tools/go/analysis/passes/waitgroup"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	mychecks := []*analysis.Analyzer{}
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}
	mychecks = append(
		mychecks,
		appends.Analyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		defers.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpmux.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		slog.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stdversion.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
		waitgroup.Analyzer,
		errcheck.Analyzer,
		ExitCheckAnalyzer,
	)
	multichecker.Main(
		mychecks...,
	)
}
