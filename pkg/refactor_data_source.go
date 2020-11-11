package pkg

import "strings"

func RefactorDataSourceType(mc *ModuleConfigs, oldType, newType, currentModuleAbsPath string) (ModuleSources, error) {
	return RefactorLabelInModule(mc, "data", "data", []string{oldType}, []string{newType}, currentModuleAbsPath,
		func(mc *ModuleConfigs, label []string) []string {
			files := []string{}
			for addr, ds := range mc.Get(currentModuleAbsPath).DataResources {
				addrs := strings.Split(addr, ".")
				if addrs[1] == label[0] {
					files = append(files, ds.DeclRange.Filename)
				}
			}
			return files
		})
}

func RefactorDataSourceName(mc *ModuleConfigs, dsType, oldName, newName, currentModuleAbsPath string) (ModuleSources, error) {
	return RefactorLabelInModule(mc, "data", "data", []string{dsType, oldName}, []string{dsType, newName}, currentModuleAbsPath,
		func(mc *ModuleConfigs, label []string) []string {
			return []string{mc.Get(currentModuleAbsPath).DataResources["data."+strings.Join(label, ".")].DeclRange.Filename}
		})
}

func RefactorDataSourceAttribute(mc *ModuleConfigs, dsType string, oldName, newName []string, currentModuleAbsPath string) (ModuleSources, error) {
	return nil, nil
}
