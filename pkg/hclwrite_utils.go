package pkg

import "github.com/hashicorp/hcl/v2/hclwrite"

func RenameVariablePrefixInBody(body *hclwrite.Body, sourceAddrs []string, destAddrs []string) {
	for _, attr := range body.Attributes() {
		attr.Expr().RenameVariablePrefix(sourceAddrs, destAddrs)
	}
	for _, blk := range body.Blocks() {
		RenameVariablePrefixInBody(blk.Body(), sourceAddrs, destAddrs)
	}
}

