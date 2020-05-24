package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestK8sOkCheck(t *testing.T) {
	targetVersion = "v1.8.0"
	analysistest.Run(t, analysistest.TestData(), Analyzer, "ok")
}

func TestK8sDeprecationCheck(t *testing.T) {
	targetVersion = "v1.15.0"
	analysistest.Run(t, analysistest.TestData(), Analyzer, "deprecation")
}

func TestK8sRemovalCheck(t *testing.T) {
	targetVersion = "v1.16.0"
	analysistest.Run(t, analysistest.TestData(), Analyzer, "removal")
}

func TestK8sMixCheck(t *testing.T) {
	targetVersion = "v1.17.0"
	analysistest.Run(t, analysistest.TestData(), Analyzer, "mix")
}
