package pkg

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/configs/configload"
	"path/filepath"
	"strings"
)

type ModuleConfigs struct {
	rootCfg           *configs.Config
	rootModuleAbsPath string
	configs           map[string]*configs.Module // key: abs path to the module; value: the module's config
}

func NewModuleConfigs( path string) (*ModuleConfigs,error) {
	loader, err := configload.NewLoader(&configload.Config{
		ModulesDir: filepath.Join(path, ".terraform", "modules"),
	})
	if err != nil {
		return nil, err
	}

	cfg, diags := loader.LoadConfig(path)
	if len(diags) != 0 {
		return nil, fmt.Errorf("failed to load config under path %q: %v", path, diags.Error())
	}
	if cfg == nil {
		return nil, errors.New("unexpected nil config")
	}
	configs := &ModuleConfigs{
		rootCfg:           cfg,
		rootModuleAbsPath: path,
		configs: map[string]*configs.Module{
			path: cfg.Module,
		},
	}

	addChildModules(cfg, path, configs.configs)
	return configs, nil
}

// addChildModules recursively adds the child modules of the provided module to fill in the moduleConfigs passed in.
func addChildModules(cfg *configs.Config, path string, moduleConfigs map[string]*configs.Module) {
	for _, childMod := range cfg.Children {
		// We skip all module sources except local paths
		if !isLocalModulePath(childMod.SourceAddr) {
			continue
		}
		childModPath := filepath.Clean(filepath.Join(path, childMod.SourceAddr))

		if _, ok := moduleConfigs[childModPath]; ok {
			continue
		}

		moduleConfigs[childModPath] = childMod.Module
		addChildModules(childMod, childModPath, moduleConfigs)
	}
}

func isLocalModulePath(path string) bool {
	// https://www.terraform.io/docs/modules/sources.html#local-paths
	return strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../")
}


func (mc ModuleConfigs) Get(modpath string) *configs.Module {
	mod := mc.configs[modpath]
	if mod == nil {
		panic(fmt.Sprintf("unexpected nil module config for path: %s", modpath))
	}
	return mod
}


