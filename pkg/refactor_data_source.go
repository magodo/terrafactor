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

func RefactorDataSourceAttribute(mc *ModuleConfigs, dsType string, dsName string, oldAddr, newAddr []string, currentModuleAbsPath string) (ModuleSources, error) {
	oldAddrs := make([]string, len(oldAddr)+2)
	newAddrs := make([]string, len(newAddr)+2)
	oldAddrs[0], newAddrs[0] = dsType, dsType
	oldAddrs[1], newAddrs[1] = dsName, dsName
	copy(oldAddrs[2:], oldAddr)
	copy(newAddrs[2:], newAddr)
	return RefactorAttributeInModule(mc, "data", "data", oldAddrs, newAddrs, currentModuleAbsPath,
		func(mc *ModuleConfigs, label []string) []string {
			return []string{mc.Get(currentModuleAbsPath).DataResources["data."+strings.Join(label, ".")].DeclRange.Filename}
		})
}
