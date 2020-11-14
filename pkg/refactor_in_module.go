package pkg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type locateDeclarationFunc func(configs *ModuleConfigs, label []string) []string

func RefactorLabelInModule(mc *ModuleConfigs, defTypeId, refTypeId string, oldLabels, newLabels []string, currentModuleAbsPath string, f locateDeclarationFunc) (ModuleSources, error) {
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
			if blk.Type() != defTypeId {
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
			file.Body().FirstMatchingBlock(defTypeId, thisLabels).SetLabels(thisNewLabels)
			moduleSources[defFile] = file.Bytes()
		}
	}

	// Load the files within this module that has referring to the object to be renamed
	for filename, src := range moduleSources {
		file, diags := hclwrite.ParseConfig(src, "", hcl.Pos{Line: 1, Column: 1})
		if len(diags) != 0 {
			return nil, fmt.Errorf("failed to parse config for %s: %v", filename, diags.Error())
		}
		oldAddr, newAddr := oldLabels, newLabels

		if refTypeId != "" {
			oldAddr = make([]string, len(oldLabels)+1)
			oldAddr[0] = refTypeId
			copy(oldAddr[1:], oldLabels)

			newAddr = make([]string, len(newLabels)+1)
			newAddr[0] = refTypeId
			copy(newAddr[1:], newLabels)
		}

		RenameVariablePrefixInBody(file.Body(), oldAddr, newAddr)

		moduleSources[filename] = file.Bytes()
	}
	return moduleSources, nil
}

func RefactorAttributeInModule(mc *ModuleConfigs, defTypeId, refTypeId string, oldAddrs, newAddrs []string, currentModuleAbsPath string, f locateDeclarationFunc) (ModuleSources, error) {
	moduleSources, err := LoadModule(currentModuleAbsPath)
	if err != nil {
		return nil, err
	}

	if len(oldAddrs) < 3 {
		return nil, fmt.Errorf("old address %q contains less than 3 segments", strings.Join(oldAddrs, "."))
	}
	if len(newAddrs) < 3 {
		return nil, fmt.Errorf("new address %q contains less than 3 segments", strings.Join(newAddrs, "."))
	}
	if len(oldAddrs) != len(newAddrs) {
		return nil, fmt.Errorf("new address %q doesn't have the same length as old address %q", strings.Join(newAddrs, "."), strings.Join(oldAddrs, "."))
	}

	oldLabels, _ := oldAddrs[:2], newAddrs[:2]
	oldAttrs, newAttrs := oldAddrs[2:], newAddrs[2:]

	// Lookup the file(s) that includes the declaration of the object(s), do the refactor on the definition
	defFiles := f(mc, oldLabels)
	for _, defFile := range defFiles {
		file, diags := hclwrite.ParseConfig(moduleSources[defFile], "", hcl.Pos{Line: 1, Column: 1})
		if len(diags) != 0 {
			return nil, errors.New(diags.Error())
		}

		topBlks := HclWriteBodyFindAllMatchingBlocks(file.Body(), defTypeId, oldLabels)

		for _, topBlk := range topBlks {
			candidateBlks := []*hclwrite.Block{topBlk}
			for i, addr := range oldAttrs[:len(oldAttrs)-1] {
				var newCandidateBlks []*hclwrite.Block
				for _, candidateBlk := range candidateBlks {
					for _, blk := range candidateBlk.Body().Blocks() {
						if blk.Type() == addr {
							newCandidateBlks = append(newCandidateBlks, blk)
						}
					}
				}
				candidateBlks = newCandidateBlks
				if len(candidateBlks) == 0 {
					return nil, fmt.Errorf("failed to find any block named %q deep in traversal %q in %s block %q", addr, strings.Join(oldAttrs[:i], "."), defTypeId, strings.Join(oldLabels, "."))
				}
			}

			// The candidateBlks contains all the blocks conforming to the oldAttrs traversal.
			// The last component of oldAttrs that is to be replaced is either an attribute or block.
			lastAttr := oldAttrs[len(oldAttrs)-1]
			newAttr := newAttrs[len(newAttrs)-1]

			for _, blk := range candidateBlks {
				if attr, ok := blk.Body().Attributes()[lastAttr]; ok {
					RenameAttributeName(attr, newAttr)
					continue
				}

				for _, blk := range blk.Body().Blocks() {
					if blk.Type() == lastAttr {
						blk.SetType(newAttr)
					}
				}
			}

			moduleSources[defFile] = file.Bytes()
		}
	}

	// Load the files within this module that has referring to the object to be renamed
	for filename, src := range moduleSources {
		file, diags := hclwrite.ParseConfig(src, "", hcl.Pos{Line: 1, Column: 1})
		if len(diags) != 0 {
			return nil, fmt.Errorf("failed to parse config for %s: %v", filename, diags.Error())
		}
		oldAddr, newAddr := oldAddrs, newAddrs

		// Managed resource will not contain the reference type.
		if refTypeId != "" {
			oldAddr = make([]string, len(oldAddrs)+1)
			oldAddr[0] = refTypeId
			copy(oldAddr[1:], oldAddrs)

			newAddr = make([]string, len(newAddrs)+1)
			newAddr[0] = refTypeId
			copy(newAddr[1:], newAddrs)
		}

		RenameVariablePrefixInBody(file.Body(), oldAddr, newAddr)

		moduleSources[filename] = file.Bytes()
	}
	return moduleSources, nil
}
