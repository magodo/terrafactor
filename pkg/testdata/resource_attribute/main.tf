
resource "foo" "a" {
  attr = "x"

  block {
    attr = 123
    nested_block {}
  }

  multi_block {
    attr1 = 1
    nested_block {}
  }

  multi_block {
    attr2 = 2
    nested_block {}
  }
}
