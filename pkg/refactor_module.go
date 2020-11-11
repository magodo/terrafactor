package pkg

func RefactorModuleName(mc *ModuleConfigs, oldName, newName, currentModuleAbsPath string) (ModuleSources, error) {
	return RefactorLabelInModule(mc, "module", "module", []string{oldName}, []string{newName}, currentModuleAbsPath,
		func(mc *ModuleConfigs, label []string) []string {
			return []string{
				mc.Get(currentModuleAbsPath).ModuleCalls[label[0]].SourceAddrRange.Filename,
			}
		})
}
