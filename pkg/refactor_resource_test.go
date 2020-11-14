package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRefactorResourceType(t *testing.T) {
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	testdataDir, _ := filepath.Abs(filepath.Join(cwd, "testdata"))
	cases := []struct {
		name           string
		rootModulePath string
		oldType        string
		newType        string
		expect         ModuleSources
	}{
		{
			name:           "basic",
			rootModulePath: filepath.Join(testdataDir, "resource_type"),
			oldType:        "type1",
			newType:        "type2",
			expect: map[string][]byte{
				filepath.Join(testdataDir, "resource_type", "main.tf"): []byte(`
resource "type2" "a" {
  name = "bar"
}
`),
				filepath.Join(testdataDir, "resource_type", "main2.tf"): []byte(`
resource "typex" "a" {
  name = type2.a.name
}
`),
			},
		},
	}

	for _, c := range cases {
		require.NoError(t, os.Chdir(c.rootModulePath), c.name)
		moduleConfigs, err := NewModuleConfigs(c.rootModulePath)
		require.NoError(t, err, c.name)
		moduleSources, err := RefactorResourceType(moduleConfigs, c.oldType, c.newType, c.rootModulePath)
		require.NoError(t, err, c.name)
		require.Equal(t, c.expect, moduleSources, c.name)
	}
}

func TestRefactorResourceName(t *testing.T) {
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	testdataDir, _ := filepath.Abs(filepath.Join(cwd, "testdata"))
	cases := []struct {
		name           string
		rootModulePath string
		resType        string
		oldName        string
		newName        string
		expect         ModuleSources
	}{
		{
			name:           "basic",
			rootModulePath: filepath.Join(testdataDir, "resource_name"),
			resType:        "foo",
			oldName:        "name1",
			newName:        "name2",
			expect: map[string][]byte{
				filepath.Join(testdataDir, "resource_name", "main.tf"): []byte(`
resource "foo" "name2" {
  name = "bar"
}
`),
				filepath.Join(testdataDir, "resource_name", "main2.tf"): []byte(`
resource "foo" "namex" {
  name = foo.name2.name
}
`),
			},
		},
	}

	for _, c := range cases {
		require.NoError(t, os.Chdir(c.rootModulePath), c.name)
		moduleConfigs, err := NewModuleConfigs(c.rootModulePath)
		require.NoError(t, err, c.name)
		moduleSources, err := RefactorResourceName(moduleConfigs, c.resType, c.oldName, c.newName, c.rootModulePath)
		require.NoError(t, err, c.name)
		require.Equal(t, c.expect, moduleSources, c.name)
	}
}

func TestRefactorReourceAttribute(t *testing.T) {
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	testdataDir, _ := filepath.Abs(filepath.Join(cwd, "testdata"))
	cases := []struct {
		name           string
		rootModulePath string
		resType        string
		resName        string
		oldAddr        []string
		newAddr        []string
		expect         ModuleSources
	}{
		{
			name:           "rename top-level attribute",
			rootModulePath: filepath.Join(testdataDir, "resource_attribute"),
			resType:        "foo",
			resName:        "a",
			oldAddr:        []string{"attr"},
			newAddr:        []string{"attr_new"},
			expect: map[string][]byte{
				filepath.Join(testdataDir, "resource_attribute", "main.tf"): []byte(`
resource "foo" "a" {
  attr_new = "x"

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
`),
				filepath.Join(testdataDir, "resource_attribute", "main2.tf"): []byte(`
resource "foo" "b" {
  attr                    = foo.a.attr_new
  attr_block              = foo.a.block.0.attr
  attr_block_nested       = foo.a.block.0.nested_block.0.attr_nest
  attr_multi_block        = foo.a.multi_block.0.attr
  attr_multi_block_nested = foo.a.multi_block.0.nested_block.attr_nest
}
`),
			},
		},
		{
			name:           "rename top-level single block",
			rootModulePath: filepath.Join(testdataDir, "resource_attribute"),
			resType:        "foo",
			resName:        "a",
			oldAddr:        []string{"block"},
			newAddr:        []string{"block_new"},
			expect: map[string][]byte{
				filepath.Join(testdataDir, "resource_attribute", "main.tf"): []byte(`
resource "foo" "a" {
  attr = "x"

  block_new {
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
`),
				filepath.Join(testdataDir, "resource_attribute", "main2.tf"): []byte(`
resource "foo" "b" {
  attr                    = foo.a.attr
  attr_block              = foo.a.block_new.0.attr
  attr_block_nested       = foo.a.block_new.0.nested_block.0.attr_nest
  attr_multi_block        = foo.a.multi_block.0.attr
  attr_multi_block_nested = foo.a.multi_block.0.nested_block.attr_nest
}
`),
			},
		},
		{
			name:           "rename top-level multi block",
			rootModulePath: filepath.Join(testdataDir, "resource_attribute"),
			resType:        "foo",
			resName:        "a",
			oldAddr:        []string{"multi_block"},
			newAddr:        []string{"multi_block_new"},
			expect: map[string][]byte{
				filepath.Join(testdataDir, "resource_attribute", "main.tf"): []byte(`
resource "foo" "a" {
  attr = "x"

  block {
    attr = 123
    nested_block {
      attr_nest = 1
    }
  }

  multi_block_new {
    attr = 1
    nested_block {
      attr_nest = 1
    }
  }

  multi_block_new {
    attr = 1
    nested_block {
      attr_nest = 1
    }
  }
}
`),
				filepath.Join(testdataDir, "resource_attribute", "main2.tf"): []byte(`
resource "foo" "b" {
  attr                    = foo.a.attr
  attr_block              = foo.a.block.0.attr
  attr_block_nested       = foo.a.block.0.nested_block.0.attr_nest
  attr_multi_block        = foo.a.multi_block_new.0.attr
  attr_multi_block_nested = foo.a.multi_block_new.0.nested_block.attr_nest
}
`),
			},
		},
		{
			name:           "rename top-level block attribute",
			rootModulePath: filepath.Join(testdataDir, "resource_attribute"),
			resType:        "foo",
			resName:        "a",
			oldAddr:        []string{"block", "attr"},
			newAddr:        []string{"block", "attr_new"},
			expect: map[string][]byte{
				filepath.Join(testdataDir, "resource_attribute", "main.tf"): []byte(`
resource "foo" "a" {
  attr = "x"

  block {
    attr_new = 123
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
`),
				filepath.Join(testdataDir, "resource_attribute", "main2.tf"): []byte(`
resource "foo" "b" {
  attr                    = foo.a.attr
  attr_block              = foo.a.block.0.attr_new
  attr_block_nested       = foo.a.block.0.nested_block.0.attr_nest
  attr_multi_block        = foo.a.multi_block.0.attr
  attr_multi_block_nested = foo.a.multi_block.0.nested_block.attr_nest
}
`),
			},
		},
		{
			name:           "rename top-level multi block nested block attr",
			rootModulePath: filepath.Join(testdataDir, "resource_attribute"),
			resType:        "foo",
			resName:        "a",
			oldAddr:        []string{"multi_block", "nested_block", "attr_nest"},
			newAddr:        []string{"multi_block", "nested_block", "attr_nest_new"},
			expect: map[string][]byte{
				filepath.Join(testdataDir, "resource_attribute", "main.tf"): []byte(`
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
      attr_nest_new = 1
    }
  }

  multi_block {
    attr = 1
    nested_block {
      attr_nest_new = 1
    }
  }
}
`),
				filepath.Join(testdataDir, "resource_attribute", "main2.tf"): []byte(`
resource "foo" "b" {
  attr                    = foo.a.attr
  attr_block              = foo.a.block.0.attr
  attr_block_nested       = foo.a.block.0.nested_block.0.attr_nest
  attr_multi_block        = foo.a.multi_block.0.attr
  attr_multi_block_nested = foo.a.multi_block.0.nested_block.attr_nest_new
}
`),
			},
		},
	}

	for _, c := range cases {
		require.NoError(t, os.Chdir(c.rootModulePath), c.name)
		moduleConfigs, err := NewModuleConfigs(c.rootModulePath)
		require.NoError(t, err, c.name)
		moduleSources, err := RefactorReourceAttribute(moduleConfigs, c.resType, c.resName, c.oldAddr, c.newAddr, c.rootModulePath)
		require.NoError(t, err, c.name)
		for sourceName := range moduleSources {
			require.Equal(t, string(c.expect[sourceName]), string(moduleSources[sourceName]), fmt.Sprintf("%s: %s", c.name, sourceName))
		}
	}
}
