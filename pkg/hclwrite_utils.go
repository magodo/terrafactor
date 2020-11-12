package pkg

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"reflect"
)

func RenameVariablePrefixInBody(body *hclwrite.Body, sourceAddrs []string, destAddrs []string) {
	for _, attr := range body.Attributes() {
		attr.Expr().RenameVariablePrefix(sourceAddrs, destAddrs)
	}
	for _, blk := range body.Blocks() {
		RenameVariablePrefixInBody(blk.Body(), sourceAddrs, destAddrs)
	}
}

func HclWriteBodyFindAllMatchingBlocks(b *hclwrite.Body, typeName string, labels []string) []*hclwrite.Block {
	var blks []*hclwrite.Block

	for _, block := range b.Blocks() {
		if typeName == block.Type() {
			labelNames := block.Labels()
			if len(labels) == 0 && len(labelNames) == 0 {
				blks = append(blks, block)
				continue
			}
			if reflect.DeepEqual(labels, labelNames) {
				blks = append(blks, block)
				continue
			}
		}
	}

	return blks
}
