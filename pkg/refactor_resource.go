package pkg

import "strings"

func RefactorResourceType(mc *ModuleConfigs, oldType, newType, currentModuleAbsPath string) (ModuleSources, error) {
	return RefactorLabelInModule(mc, "resource", "", []string{oldType}, []string{newType}, currentModuleAbsPath,
		func(mc *ModuleConfigs, label []string) []string {
			files := []string{}
			for addr, ds := range mc.Get(currentModuleAbsPath).ManagedResources {
				addrs := strings.Split(addr, ".")
				if addrs[0] == label[0] {
					files = append(files, ds.DeclRange.Filename)
				}
			}
			return files
		})
}

func RefactorResourceName(mc *ModuleConfigs, resType, oldName, newName, currentModuleAbsPath string) (ModuleSources, error) {
	return RefactorLabelInModule(mc, "resource", "", []string{resType, oldName}, []string{resType, newName}, currentModuleAbsPath,
		func(mc *ModuleConfigs, label []string) []string {
			return []string{mc.Get(currentModuleAbsPath).ManagedResources[strings.Join(label, ".")].DeclRange.Filename}
		})
}
