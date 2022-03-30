package lint

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOsExitForbiddenAnalyzer(t *testing.T) {

	dir, cleanup, err := analysistest.WriteFiles(map[string]string{"sample/osexitforbidden.go": `
package main

import os1 "os"

func main() {
	os1.Exit(1) // want "os.Exit call is forbidden in main func of main package"
}`,
	})

	defer cleanup()
	require.NoError(t, err, "samples write failed")

	analysistest.Run(t, dir, OsExitForbiddenAnalyzer, "sample")
}
