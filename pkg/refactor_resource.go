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

func RefactorReourceAttribute(mc *ModuleConfigs, resType string, resName string, oldAddr, newAddr []string, currentModuleAbsPath string) (ModuleSources, error) {
	oldAddrs := make([]string, len(oldAddr)+2)
	newAddrs := make([]string, len(newAddr)+2)
	oldAddrs[0], newAddrs[0] = resType, resType
	oldAddrs[1], newAddrs[1] = resName, resName
	copy(oldAddrs[2:], oldAddr)
	copy(newAddrs[2:], newAddr)
	return RefactorAttributeInModule(mc, "resource", "", oldAddrs, newAddrs, currentModuleAbsPath,
		func(mc *ModuleConfigs, label []string) []string {
			return []string{mc.Get(currentModuleAbsPath).ManagedResources[strings.Join(label, ".")].DeclRange.Filename}
		})
}
