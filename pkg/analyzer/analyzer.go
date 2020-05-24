package analyzer

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/mod/semver"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	//"golang.org/x/tools/go/analysis/unitchecker"

	"github.com/yoichiwo7/k8sdepr/pkg/version"
)

const Doc = `find kubernetes api deprecation/removal in target version.
-targetVersion must be set to run the check.`

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

type kubeAPISelector struct {
	APIVersion string
	Kind       string
	Selector   *ast.SelectorExpr
}

func run(pass *analysis.Pass) (interface{}, error) {
	if !semver.IsValid(targetVersion) {
		return nil, fmt.Errorf("invalid semver: %v", targetVersion)
	}

	deprecationMap := version.GetDeprecationMap()

	for _, f := range pass.Files {
		importPathMap := getImportPathMap(f.Imports)
		// Get import map (key=name, value=path)

		ast.Inspect(f, func(n ast.Node) bool {
			// Parse package and kind struct
			k8sSelector, err := parseKubeAPISelector(n, importPathMap)
			if err != nil {
				return true
			}

			// Check Kubernetes API deprecation/removal
			deprInfo, ok := deprecationMap[version.DeprecationKey{APIVersion: k8sSelector.APIVersion, Kind: k8sSelector.Kind}]
			if !ok {
				return true
			}
			if k8sSelector.APIVersion != deprInfo.APIVersion || k8sSelector.Kind != deprInfo.Kind {
				return true
			}
			if !ignoreRemoval && semver.Compare(targetVersion, deprInfo.RemovedIn) >= 0 {
				pass.Report(analysis.Diagnostic{
					Pos: k8sSelector.Selector.Pos(),
					Message: fmt.Sprintf("%v:%v is removed in %v. Migrate to %v:%v.",
						k8sSelector.APIVersion, k8sSelector.Kind, deprInfo.RemovedIn, deprInfo.ReplacementAPI, k8sSelector.Kind),
				})
				return true
			}
			if !ignoreDeprecation && semver.Compare(targetVersion, deprInfo.DeprecatedIn) >= 0 {
				pass.Report(analysis.Diagnostic{
					Pos: k8sSelector.Selector.Pos(),
					Message: fmt.Sprintf("%v:%v is deprecated in %v. Migrate to %v:%v.",
						k8sSelector.APIVersion, k8sSelector.Kind, deprInfo.DeprecatedIn, deprInfo.ReplacementAPI, k8sSelector.Kind),
				})
				return true
			}
			return true
		})
	}

	return nil, nil
}

func parseKubeAPISelector(n ast.Node, importPathMap map[string]string) (*kubeAPISelector, error) {
	retErr := fmt.Errorf("failed to parse selector")

	selector, ok := n.(*ast.SelectorExpr)
	if !ok {
		return nil, retErr
	}
	xIdent, ok := selector.X.(*ast.Ident)
	if !ok {
		return nil, retErr
	}
	kind := selector.Sel.Name
	pkgName := xIdent.Name

	// Get matched apiVersion
	apiVersion, ok := importPathMap[pkgName]
	if !ok {
		return nil, retErr
	}

	return &kubeAPISelector{
		APIVersion: apiVersion,
		Kind:       kind,
		Selector:   selector,
	}, nil

}

func getImportPathMap(imports []*ast.ImportSpec) map[string]string {
	importMap := map[string]string{}
	for _, imp := range imports {
		path := trimPath(imp.Path.Value)
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
	return importMap
}

func trimPath(path string) string {
	return strings.ReplaceAll(path, "\"", "")
}

// foundFact is a fact associated with functions that match -name.
// We use it to exercise the fact machinery in tests.
type foundFact struct{}

func (*foundFact) String() string { return "found" }
func (*foundFact) AFact()         {}
