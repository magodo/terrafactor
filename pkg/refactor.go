package pkg

import (
	"io/ioutil"
	"os"
)

func Refactor(rootModuleAbsPath string, sourceAddr []string, destAddr []string, currentModuleAbsPath string) error {
	// We have to chdir here because otherwise the loader.LoadConfig will fail during looking up
	// the module based on its manifest record, where the path is recorded in a relative manner.
	if err := os.Chdir(rootModuleAbsPath); err != nil {
		return err
	}
	moduleConfigs, err := NewModuleConfigs(rootModuleAbsPath)
	if err != nil {
		return err
	}

	var moduleSources ModuleSources

	switch sourceAddr[0] {
	case "module":
		switch len(sourceAddr) {
		case 2:
			moduleSources, err = RefactorModuleName(moduleConfigs, sourceAddr[1], destAddr[1], currentModuleAbsPath)
			if err != nil {
				return err
			}
		default:
			panic("TODO: rename the module output variable attribute (cross module)")
		}
	case "data":
		switch len(sourceAddr) {
		case 2:
			moduleSources, err = RefactorDataSourceType(moduleConfigs, sourceAddr[1], destAddr[1], currentModuleAbsPath)
		case 3:
			moduleSources, err = RefactorDataSourceName(moduleConfigs, sourceAddr[1], sourceAddr[2], destAddr[2], currentModuleAbsPath)
		default:
			moduleSources, err = RefactorDataSourceAttribute(moduleConfigs, sourceAddr[1], sourceAddr[2], sourceAddr[3:], destAddr[3:], currentModuleAbsPath)
		}
	case "var":
		switch len(sourceAddr) {
		case 2:
			panic("TODO: rename the input variable name (cross module)")
		default:
			panic("TODO: rename the input variable attribute (for complex type) (cross module)")
		}
	case "output":
		switch len(sourceAddr) {
		case 2:
			panic("TODO: rename the output variable name (cross module)")
		default:
			panic("TODO: rename the output variable attribute (for complex type) (cross module)")
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
		case 1:
			moduleSources, err = RefactorResourceType(moduleConfigs, sourceAddr[0], destAddr[0], currentModuleAbsPath)
		case 2:
			moduleSources, err = RefactorResourceName(moduleConfigs, sourceAddr[0], sourceAddr[1], destAddr[1], currentModuleAbsPath)
		default:
			moduleSources, err = RefactorReourceAttribute(moduleConfigs, sourceAddr[0], sourceAddr[1], sourceAddr[2:], destAddr[2:], currentModuleAbsPath)
		}
	}

	// write back the module sources
	for f, b := range moduleSources {
		if err := ioutil.WriteFile(f, b, 0644); err != nil {
			return err
		}
	}
	return nil
}
