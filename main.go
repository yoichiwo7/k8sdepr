package main

import (
	"github.com/yoichiwo7/k8sdepr/pkg/analyzer"
	sg "golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	sg.Main(
		analyzer.Analyzer,
	)
}
