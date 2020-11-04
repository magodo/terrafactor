package pkg

import (
	"io/ioutil"
	"os"
)

func Refactor(rootModuleAbsPath string,  sourceAddr []string, destAddr []string, currentModuleAbsPath string) error {
	// We have to chdir here because otherwise the loader.LoadConfig will fail during looking up
	// the module based on its manifest record, where the path is recorded in a relative manner.
	if err := os.Chdir(rootModuleAbsPath); err != nil {
		return err
	}
	moduleConfigs,err := NewModuleConfigs(rootModuleAbsPath)
	if err != nil {
		return err
	}

	var moduleSources ModuleSources

	switch sourceAddr[0] {
	case "module":
		switch len(sourceAddr) {
		case 2:
			moduleSources, err = RefactorModuleName(moduleConfigs, sourceAddr, destAddr[0], currentModuleAbsPath)
			if err != nil {
				return err
			}
		default:
			panic("TODO: rename the module output variable attribute")
		}
	case "data":
		switch len(sourceAddr) {
		case 2:
			panic("TODO: rename the data source type")
		default:
			panic("TODO: rename the data source attribute")
		}
	case "var":
		switch len(sourceAddr) {
		case 2:
			panic("TODO: rename the input variable name")
		default:
			panic("TODO: rename the input variable attribute (for complex type)")
		}
	case "local":
		switch len(sourceAddr) {
		case 2:
			panic("TODO: rename the local variable name")
		default:
			panic("TODO: rename the local variable attribute (for complex type)")
		}
	default:
		switch len(sourceAddr) {
		case 2:
			panic("TODO: rename the resource type")
		default:
			panic("TODO: rename the resource attribute")
		}
	}

	// write back the module sources
	for f, b := range moduleSources {
		if err := ioutil.WriteFile(f, b, 0644); err  != nil {
			return err
		}
	}
	return nil
}
