package analyzer

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/mod/semver"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	//"golang.org/x/tools/go/analysis/unitchecker"
)

const Doc = `find kubernetes api deprecation/removal in target version.`

const KubeAPIImportPrefix = "k8s.io/api/"
const KubeCoreAPIPrefix = "core/"

var Analyzer = &analysis.Analyzer{
	Name:             "k8sdepr",
	Doc:              Doc,
	Run:              run,
	RunDespiteErrors: true,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	FactTypes:        []analysis.Fact{new(foundFact)},
}

var targetVersion string   // -targetVersion flag
var ignoreDeprecation bool // -ignoreDeprecation flag
var ignoreRemoval bool     // -ignoreRemoval flag

func init() {
	// NOTE: If you want to pass parameter from 'go vet -vettool=...' you need to use flag module
	Analyzer.Flags.StringVar(&targetVersion, "targetVersion", targetVersion, "target semantic version of the Kubernetes (ex. v1.16.0)")
	Analyzer.Flags.BoolVar(&ignoreDeprecation, "ignoreDeprecation", ignoreDeprecation, "ignore deprecation detection")
	Analyzer.Flags.BoolVar(&ignoreRemoval, "ignoreRemoval", ignoreRemoval, "ignore removal detection")
}

func run(pass *analysis.Pass) (interface{}, error) {
	if !semver.IsValid(targetVersion) {
		return nil, fmt.Errorf("invalid semver: %v", targetVersion)
	}

	for _, f := range pass.Files {
		// Get import map (key=name, value=path)
		importMap := map[string]string{}
		for _, imp := range f.Imports {
			path := strings.ReplaceAll(imp.Path.Value, "\"", "")
			localPkgName := imp.Name
			if !strings.HasPrefix(path, KubeAPIImportPrefix) {
				continue
			}
			apiVersion := strings.SplitAfter(path, KubeAPIImportPrefix)[1]
			if strings.HasPrefix(apiVersion, KubeCoreAPIPrefix) {
				apiVersion = strings.SplitAfter(apiVersion, KubeCoreAPIPrefix)[1]
			}

			if localPkgName != nil {
				// has local import name
				importMap[localPkgName.String()] = apiVersion
			} else {
				elems := strings.Split(path, "/")
				pkgName := elems[len(elems)-1]
				importMap[pkgName] = apiVersion
			}
		}

		ast.Inspect(f, func(n ast.Node) bool {
			// Parse package and kind struct
			selector, ok := n.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			xIdent, ok := selector.X.(*ast.Ident)
			if !ok {
				return true
			}
			kind := selector.Sel.Name
			pkgName := xIdent.Name

			// Get matched apiVersion
			apiVersion, ok := importMap[pkgName]
			if !ok {
				return true
			}

			// Check deprecated/removed API
			// TODO: Use map for better performance
			for _, di := range DeprecationInfoList {
				if apiVersion != di.Name || kind != di.Kind {
					continue
				}
				if !ignoreRemoval && semver.Compare(targetVersion, di.RemovedIn) >= 0 {
					pass.Report(analysis.Diagnostic{
						Pos: selector.Pos(),
						Message: fmt.Sprintf("%v:%v is removed in %v. Migrate to %v:%v.",
							apiVersion, kind, di.RemovedIn, di.ReplacementAPI, kind),
					})
					break
				}
				if !ignoreDeprecation && semver.Compare(targetVersion, di.DeprecatedIn) >= 0 {
					pass.Report(analysis.Diagnostic{
						Pos: selector.Pos(),
						Message: fmt.Sprintf("%v:%v is deprecated in %v. Migrate to %v:%v.",
							apiVersion, kind, di.DeprecatedIn, di.ReplacementAPI, kind),
					})
					break
				}
			}
			return true
		})
	}

	return nil, nil
}

// foundFact is a fact associated with functions that match -name.
// We use it to exercise the fact machinery in tests.
type foundFact struct{}

func (*foundFact) String() string { return "found" }
func (*foundFact) AFact()         {}
