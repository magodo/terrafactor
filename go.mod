module github.com/magodo/terrafactor

go 1.15

require (
	github.com/hashicorp/hcl/v2 v2.7.0
	github.com/hashicorp/terraform v0.13.5
	github.com/stretchr/testify v1.5.1
)

replace github.com/hashicorp/hcl/v2 => github.com/magodo/hcl/v2 v2.3.1-0.20201115055447-6fbb048a5c3a
