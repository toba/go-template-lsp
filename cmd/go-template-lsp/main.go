// Command go-template-lsp provides a Language Server Protocol server for Go templates.
//
// Based on https://github.com/yayolande/go-template-lsp (MIT License)
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	tmpl "github.com/toba/go-template-lsp/internal/template"
	"github.com/toba/lsp/server"
)

// version is set by goreleaser at build time.
var version = "dev"

// TargetFileExtensions lists the file extensions this LSP supports.
var TargetFileExtensions = []string{
	"go.html", "go.tmpl", "go.txt",
	"gohtml", "gotmpl", "tmpl", "tpl",
	"html",
}

const serverName = "Gozer"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "check" {
		checkFlags := flag.NewFlagSet("check", flag.ExitOnError)
		jsonFlag := checkFlags.Bool("json", false, "output as JSON")
		funcsFlag := checkFlags.String(
			"funcs",
			"",
			"comma-separated list of additional template function names",
		)
		_ = checkFlags.Parse(os.Args[2:])
		var extraFuncs []string
		if *funcsFlag != "" {
			extraFuncs = strings.Split(*funcsFlag, ",")
		}
		os.Exit(runCheck(checkFlags.Args(), *jsonFlag, extraFuncs))
	}

	versionFlag := flag.Bool("version", false, "print the LSP version")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("%s -- version %s\n", serverName, version)
		os.Exit(0)
	}

	srv := &server.Server{
		Name:    serverName,
		Version: version,
		Handler: newHandler(),
	}

	if err := srv.Run(context.Background()); err != nil {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}

// customFunctionsEqual returns true if two custom function maps have the same keys.
func customFunctionsEqual(
	a, b map[string]*tmpl.FunctionDefinition,
) bool {
	if len(a) != len(b) {
		return false
	}
	for name := range a {
		if _, ok := b[name]; !ok {
			return false
		}
	}
	return true
}
