package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/magodo/terrafactor/pkg"
)

var (
	flagRootModulePath    *string
	flagCurrentModulePath *string
	flagSourceAddr        *string
	flagDestAttr          *string
)

func init() {
	flagRootModulePath = flag.String("root-module-path", "", "path to root module")
	flagCurrentModulePath = flag.String("module-path", "", "path to the current module")
	flagSourceAddr = flag.String("source-addr", "", "the source attribute address resides in the current module that is to refactor")
	flagDestAttr = flag.String("dest-attr", "", "the destination attribute to rename to")
}

// convertTraversal converts the traversal into string slice.
// It will also checks whether the traversal contains any index or splat traverse, which is not
// supported by the RenameVariablePrefix in hclwrite for now.
func convertTraversal(traversal hcl.Traversal) ([]string, error) {
	out := []string{}
	for _, traverse := range traversal {
		switch traverse := traverse.(type) {
		case hcl.TraverseIndex:
			return nil, fmt.Errorf("index traverse is not supported")
		case hcl.TraverseSplat:
			return nil, fmt.Errorf("splat traverse is not supported")
		case hcl.TraverseRoot:
			out = append(out, traverse.Name)
		case hcl.TraverseAttr:
			out = append(out, traverse.Name)
		}
	}
	return out, nil
}

func main() {
	flag.Parse()

	// resolve the paths before we chdir
	rootModuleAbsPath, err := filepath.Abs(*flagRootModulePath)
	if err != nil {
		log.Fatal(err)
	}
	currentModuleAbsPath, err := filepath.Abs(*flagCurrentModulePath)
	if err != nil {
		log.Fatal(err)
	}

	sourceTraversal, diags := hclsyntax.ParseTraversalAbs([]byte(*flagSourceAddr), "", hcl.Pos{Line: 1, Column: 1})
	if len(diags) != 0 {
		log.Fatal(diags.Error())
	}

	sourceAddr, err := convertTraversal(sourceTraversal)
	if err != nil {
		log.Fatal(err)
	}
	if len(sourceAddr) < 2 {
		log.Fatalf("invalid source address after simplification: %s", *flagSourceAddr)
	}

	destTraversal, diags := hclsyntax.ParseTraversalAbs([]byte(*flagDestAttr), "", hcl.Pos{Line: 1, Column: 1})
	if len(diags) != 0 {
		log.Fatal(diags.Error())
	}

	destAttrs, err := convertTraversal(destTraversal)
	if err != nil {
		log.Fatal(err)
	}
	if len(destAttrs) != 1 {
		log.Fatalf("the destination address %q should only contain the one segment (the last segment)", *flagDestAttr)
	}
	destAddr := make([]string, len(sourceAddr))
	copy(destAddr, sourceAddr)
	destAddr[len(destAddr)-1] = destAttrs[0]

	if err := pkg.Refactor(rootModuleAbsPath, sourceAddr, destAddr, currentModuleAbsPath); err != nil {
		log.Fatal(err)
	}
}
