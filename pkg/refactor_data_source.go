package pkg

import "strings"

func RefactorDataSourceType(mc *ModuleConfigs, oldType, newType, currentModuleAbsPath string) (ModuleSources, error) {
	return RefactorLabelInModule(mc, "data", []string{oldType}, []string{newType}, currentModuleAbsPath,
		func(mc *ModuleConfigs, label []string ) []string {
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
