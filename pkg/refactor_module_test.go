package pkg

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestRefactorModuleName(t *testing.T) {
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	testdataDir, _ := filepath.Abs(filepath.Join(cwd, "testdata"))
	cases := []struct{
		name           string
		rootModulePath string
		oldName        string
		newName        string
		expect         ModuleSources
	}{
		{
			name:           "basic",
			rootModulePath: filepath.Join(testdataDir, "module_name"),
			oldName:        "mod1",
			newName:        "mod2",
			expect: map[string][]byte{
				filepath.Join( testdataDir, "module_name", "main.tf"): []byte(`
resource "foo" "label" {
  a = module.mod2.x
}

module "mod2" {
  source = "./mod1"
}
`),
			},
		},
	}

	for _, c := range cases {
		require.NoError(t, os.Chdir(c.rootModulePath), c.name)
		moduleConfigs,err := NewModuleConfigs(c.rootModulePath)
		require.NoError(t, err, c.name)
		moduleSources, err := RefactorModuleName(moduleConfigs,  c.oldName,  c.newName, c.rootModulePath)
		require.NoError(t, err, c.name)
		require.Equal(t, c.expect, moduleSources, c.name)
	}
}
