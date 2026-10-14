// Package modules is the single list of optional modules compiled into the app.
//
// Add a module: create internal/modules/<name> implementing module.Module, then
// append one Entry below. Remove a module: delete its folder and its Entry (and
// its tests under test/). MODULES (comma-separated) picks which entries run;
// when unset, every entry with Default: true runs.
package modules

import (
	"fmt"
	"strings"

	"be/internal/config"
	"be/internal/module"
	"be/internal/modules/webhooks"
)

// Entry registers one module.
type Entry struct {
	Name    string
	Default bool
	New     func(*module.Kernel) module.Module
}

// Available lists every module compiled into this binary, in wiring order.
var Available = []Entry{
	{Name: webhooks.Name, Default: true, New: webhooks.New},
}

// Enabled resolves which entries run for cfg (MODULES env or registry defaults).
func Enabled(cfg config.Config) ([]Entry, error) {
	if !cfg.ModulesSet {
		var out []Entry
		for _, e := range Available {
			if e.Default {
				out = append(out, e)
			}
		}
		return out, nil
	}

	byName := make(map[string]Entry, len(Available))
	for _, e := range Available {
		byName[e.Name] = e
	}

	want := make(map[string]bool, len(cfg.Modules))
	for _, name := range cfg.Modules {
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "" {
			continue
		}
		if _, ok := byName[name]; !ok {
			return nil, fmt.Errorf("modules: unknown module %q in MODULES (available: %s)", name, availableNames())
		}
		want[name] = true
	}

	var out []Entry
	for _, e := range Available { // keep registry order
		if want[e.Name] {
			out = append(out, e)
		}
	}
	return out, nil
}

// Build constructs the enabled modules against the kernel.
func Build(k *module.Kernel) (module.Set, error) {
	entries, err := Enabled(k.Config)
	if err != nil {
		return nil, err
	}
	set := make(module.Set, 0, len(entries))
	for _, e := range entries {
		set = append(set, e.New(k))
	}
	return set, nil
}

func availableNames() string {
	names := make([]string, 0, len(Available))
	for _, e := range Available {
		names = append(names, e.Name)
	}
	return strings.Join(names, ", ")
}
