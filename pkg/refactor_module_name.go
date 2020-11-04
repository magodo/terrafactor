package pkg

import (
	"errors"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func RefactorModuleName(mc *ModuleConfigs, sourceAddrs []string, destAddr string, currentModuleAbsPath string) (ModuleSources, error) {
	moduleSources, err := LoadModule(currentModuleAbsPath)
	if err != nil {
		return nil, err
	}

	moduleType, oldModuleName, newModuleName := sourceAddrs[0], sourceAddrs[1], destAddr

	// Load the file that includes the module definition, do the renaming
	modDefFile := mc.Get(currentModuleAbsPath).ModuleCalls[oldModuleName].SourceAddrRange.Filename
	file, diags := hclwrite.ParseConfig(moduleSources[modDefFile], "", hcl.Pos{Line: 1, Column: 1})
	if len(diags) != 0 {
		return nil, errors.New(diags.Error())
	}
	file.Body().FirstMatchingBlock(moduleType, []string{oldModuleName}).SetLabels([]string{newModuleName})

	moduleSources[modDefFile] = file.Bytes()

	// Load the files within this module that has referring to the module to be renamed
	for filename, src := range moduleSources {
		file, diags := hclwrite.ParseConfig(src, "", hcl.Pos{Line: 1, Column: 1})
		if len(diags) != 0 {
			return nil, fmt.Errorf("failed to parse config for %s: %v", filename, diags.Error())
		}
		destAddrs := make([]string, len(sourceAddrs))
		copy(destAddrs, sourceAddrs)
		destAddrs[len(destAddrs)-1] = destAddr
		RenameVariablePrefixInBody(file.Body(), sourceAddrs, destAddrs)

		moduleSources[filename] = file.Bytes()
	}
	return moduleSources, nil
}

