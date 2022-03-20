package main

import (
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/zhupanovdm/go-runtime-monitor/pkg/lint"
)

func main() {
	var checks []*analysis.Analyzer

	pattern := regexp.MustCompile(`SA\d+`)
	for _, v := range staticcheck.Analyzers {
		if pattern.MatchString(v.Name) {
			checks = append(checks, v)
		}
	}

	checks = append(checks, lint.StandardAnalyzers()...)
	checks = append(checks, stylecheck.Analyzers["ST1003"])
	checks = append(checks, simple.Analyzers["S1007"])
	checks = append(checks, lint.OsExitForbiddenAnalyzer)

	multichecker.Main(checks...)
}
