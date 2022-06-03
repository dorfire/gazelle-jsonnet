package jsonnet

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func (*jsonnetLang) Name() string {
	return jsonnetName
}

func (*jsonnetLang) Resolve(c *config.Config, ix *resolve.RuleIndex, _ *repo.RemoteCache, r *rule.Rule, imports interface{}, from label.Label) {
	if !shouldProcessPkg(from.Pkg) {
		log.Printf("jsonnet: Resolve() skipped for rule %s in pkg %s", r.Name(), from.Pkg)
		return
	}

	libsonnetImports, ok := imports.([]string)
	if !ok {
		log.Printf("jsonnet: Resolve() skipped for invalid `imports` type %T", imports)
		return
	}

	var bazelDeps []string
	for _, lib := range libsonnetImports {
		if !strings.HasSuffix(lib, jsonnetLibExt) {
			// Hack to prevent logs on ".json" imports which should be ignored by gazelle.
			if !strings.HasSuffix(lib, ".json") {
				log.Printf("jsonnet: ignoring unsupported import '%s' in '%s'", lib, from.Pkg)
			}
			continue
		}

		imp := libsonnetTarget(path.Split(lib))
		_, err := findIndexedRule(imp, ix, c)
		if err != nil {
			log.Printf("jsonnet: could not resolve import from '%s': %s", from, err)
			continue
		}
		bazelDeps = append(bazelDeps, imp)
	}

	if len(bazelDeps) == 0 {
		r.DelAttr("deps")
		return
	}

	r.SetAttr("deps", bazelDeps)
}

func (l *jsonnetLang) Imports(_ *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	// Generated jsonnet_to_json rules don't need to be indexed (since nothing imports them)
	if !shouldProcessPkg(f.Pkg) || r.Kind() != "jsonnet_library" {
		return nil
	}

	// Index generated jsonnet_library rules for later use in Resolve()
	return []resolve.ImportSpec{{Lang: jsonnetName, Imp: libsonnetTarget(f.Pkg, r.Name())}}
}

func (*jsonnetLang) Embeds(*rule.Rule, label.Label) []label.Label {
	return []label.Label{}
}

func findIndexedRule(target string, ix *resolve.RuleIndex, c *config.Config) (resolve.FindResult, error) {
	rules := ix.FindRulesByImportWithConfig(c, resolve.ImportSpec{Lang: jsonnetName, Imp: target}, jsonnetName)
	switch len(rules) {
	case 0:
		return resolve.FindResult{}, fmt.Errorf("no rule found for import '%s'", target)
	case 1:
		return rules[0], nil
	default:
		return resolve.FindResult{}, fmt.Errorf("multiple rules found for import '%s'", target)
	}
}

// libsonnetTarget returns the Bazel target path for the given Jsonnet library (.libsonnet file).
// e.g. 'deployments/k8s', 'consts.libsonnet' => '//deployments/k8s:consts'
func libsonnetTarget(pkg, libName string) string {
	dirPath := strings.TrimSuffix(pkg, "/")
	fileName := strings.TrimSuffix(libName, jsonnetLibExt)
	return fmt.Sprintf("//%s:%s", dirPath, fileName)
}
