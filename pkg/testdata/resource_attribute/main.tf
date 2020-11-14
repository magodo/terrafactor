
resource "foo" "a" {
  attr = "x"

  block {
    attr = 123
    nested_block {
      attr_nest = 1
    }
  }

  multi_block {
    attr = 1
    nested_block {
      attr_nest = 1
    }
  }

  multi_block {
    attr = 1
    nested_block {
      attr_nest = 1
    }
  }
}
