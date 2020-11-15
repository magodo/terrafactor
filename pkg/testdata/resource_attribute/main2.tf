
resource "bar" "c" {
  attr                    = foo.a.attr
  attr_block              = foo.a.block.0.attr
  attr_block_nested       = foo.a.block.0.nested_block.0.attr_nest
  attr_multi_block        = foo.a.multi_block.0.attr
  attr_multi_block_nested = foo.a.multi_block.0.nested_block.attr_nest
  attr2                   = foo.b.attr
  attr_block2             = foo.b.block.0.attr
}
