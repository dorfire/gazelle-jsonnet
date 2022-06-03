package jsonnet

import (
	"flag"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func (*jsonnetLang) KnownDirectives() []string {
	return []string{"jsonnet_to_json", "jsonnet_library"}
}

// RegisterFlags is required to satisfy the interface language.Language. Unused by this language.
func (*jsonnetLang) RegisterFlags(_ *flag.FlagSet, _ string, _ *config.Config) {
}

// Configure modifies the configuration using directives and other information
// extracted from a build file. Configure is called in each directory.
// Required to satisfy the interface language.Language. Unused by this language.
func (*jsonnetLang) Configure(_ *config.Config, _ string, _ *rule.File) {
}

// CheckFlags is required to satisfy the interface language.Language. Unused by this language.
func (*jsonnetLang) CheckFlags(_ *flag.FlagSet, _ *config.Config) error {
	return nil
}

// shouldProcessPkg decides if the extension should run for a given Bazel package. TODO: inspect the configuration
func shouldProcessPkg(pkg string) bool {
	return strings.HasPrefix(pkg, "cmd/") || strings.HasPrefix(pkg, "deployments/") || strings.HasPrefix(pkg, "third_party/jsonnet/")
}
