package pkg

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestRefactorModuleName(t *testing.T) {
	cwd, _ := os.Getwd()
	testdataDir, _ := filepath.Abs(filepath.Join(cwd, "..", "testdata"))
	cases := []struct{
		name string
		rootModulePath string
		oldModuleName string
		newModuleName string
		expect ModuleSources
	}{
		{
			name: "basic",
			rootModulePath: filepath.Join(testdataDir, "module_name"),
			oldModuleName:  "mod1",
			newModuleName:  "mod2",
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
		moduleSources, err := RefactorModuleName(moduleConfigs, []string{"module", c.oldModuleName}, c.newModuleName, c.rootModulePath)
		require.NoError(t, err, c.name)
		require.Equal(t, c.expect, moduleSources, c.name)
	}
}
