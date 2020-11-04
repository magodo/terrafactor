
resource "foo" "label" {
  a = module.mod1.x
}

module "mod1" {
  source = "./mod1"
}
