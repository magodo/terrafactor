package pkg

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func RenameVariablePrefixInBody(body *hclwrite.Body, sourceAddrs []string, destAddrs []string) {
	for _, attr := range body.Attributes() {
		attr.Expr().RenameVariablePrefix(sourceAddrs, destAddrs)
	}
	for _, blk := range body.Blocks() {
		RenameVariablePrefixInBody(blk.Body(), sourceAddrs, destAddrs)
	}
}

func RenameAttributeName(attr *hclwrite.Attribute, name string) {
	attr.BuildTokens(nil)[0].Bytes = []byte(name)
}

// HclWriteBodyFindAllMatchingBlocks checks the top level blocks of the body and returns those matching the type name
// and labels (if any). Each label component could be "*" to match any.
func HclWriteBodyFindAllMatchingBlocks(b *hclwrite.Body, typeName string, labels []string) []*hclwrite.Block {
	var blks []*hclwrite.Block

BlockLoop:
	for _, block := range b.Blocks() {
		if typeName == block.Type() {
			labelNames := block.Labels()
			if len(labels) != len(labelNames) {
				continue
			}
			for i := range labels {
				if labels[i] != labelNames[i] && labels[i] != "*" {
					continue BlockLoop
				}
			}
			blks = append(blks, block)
		}
	}

	return blks
}
