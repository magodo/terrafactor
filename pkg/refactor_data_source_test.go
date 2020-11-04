package pkg

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestRefactorDataSourceType(t *testing.T) {
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	testdataDir, _ := filepath.Abs(filepath.Join(cwd, "testdata"))
	cases := []struct{
		name           string
		rootModulePath string
		oldType        string
		newType        string
		expect         ModuleSources
	}{
		{
			name:           "basic",
			rootModulePath: filepath.Join(testdataDir, "data_source_name"),
			oldType:        "type1",
			newType:        "type2",
			expect: map[string][]byte{
				filepath.Join( testdataDir, "data_source_name", "main.tf"): []byte(`
data "type2" "name" {
  name = "bar"
}

module "mod1" {
  source = "./mod1"
}
`),
				filepath.Join( testdataDir, "data_source_name", "main2.tf"): []byte(`
data "typex" "name" {
  name = data.type2.name
}
`),
			},
		},
	}

	for _, c := range cases {
		require.NoError(t, os.Chdir(c.rootModulePath), c.name)
		moduleConfigs,err := NewModuleConfigs(c.rootModulePath)
		require.NoError(t, err, c.name)
		moduleSources, err := RefactorDataSourceType(moduleConfigs,  c.oldType,  c.newType, c.rootModulePath)
		require.NoError(t, err, c.name)
		require.Equal(t, c.expect, moduleSources, c.name)
	}
}
