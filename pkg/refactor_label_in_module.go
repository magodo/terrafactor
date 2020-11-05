package pkg

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type locateDeclarationFunc func(configs *ModuleConfigs, label []string) []string

func RefactorLabelInModule(mc *ModuleConfigs, blockType string, oldLabels, newLabels []string, currentModuleAbsPath string, f locateDeclarationFunc) (ModuleSources, error) {
	moduleSources, err := LoadModule(currentModuleAbsPath)
	if err != nil {
		return nil, err
	}

	// Lookup the file(s) that includes the declaration of the object(s), do the refactor on the definition
	defFiles := f(mc, oldLabels)
	for _, defFile := range defFiles {
		file, diags := hclwrite.ParseConfig(moduleSources[defFile], "", hcl.Pos{Line: 1, Column: 1})
		if len(diags) != 0 {
			return nil, errors.New(diags.Error())
		}

	blockLoop:
		for _, blk := range file.Body().Blocks() {
			if blk.Type() != blockType {
				continue
			}
			thisLabels := blk.Labels()
			for idx := range oldLabels {
				if thisLabels[idx] != oldLabels[idx] {
					continue blockLoop
				}
			}
			thisNewLabels := make([]string, len(thisLabels))
			copy(thisNewLabels, thisLabels)
			copy(thisNewLabels, newLabels)
			file.Body().FirstMatchingBlock(blockType, thisLabels).SetLabels(thisNewLabels)
			moduleSources[defFile] = file.Bytes()
		}
	}

	// Load the files within this module that has referring to the object to be renamed
	for filename, src := range moduleSources {
		file, diags := hclwrite.ParseConfig(src, "", hcl.Pos{Line: 1, Column: 1})
		if len(diags) != 0 {
			return nil, fmt.Errorf("failed to parse config for %s: %v", filename, diags.Error())
		}
		oldAddr := make([]string, len(oldLabels)+1)
		oldAddr[0] = blockType
		copy(oldAddr[1:], oldLabels)

		newAddr := make([]string, len(newLabels)+1)
		newAddr[0] = blockType
		copy(newAddr[1:], newLabels)

		RenameVariablePrefixInBody(file.Body(), oldAddr, newAddr)

		moduleSources[filename] = file.Bytes()
	}
	return moduleSources, nil
}
