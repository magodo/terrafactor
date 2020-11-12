package pkg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRefactorDataSourceType(t *testing.T) {
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
			rootModulePath: filepath.Join(testdataDir, "data_source_type"),
			oldType:        "type1",
			newType:        "type2",
			expect: map[string][]byte{
				filepath.Join(testdataDir, "data_source_type", "main.tf"): []byte(`
data "type2" "a" {
  name = "bar"
}
`),
				filepath.Join(testdataDir, "data_source_type", "main2.tf"): []byte(`
data "typex" "a" {
  name = data.type2.a.name
}
`),
			},
		},
	}

	for _, c := range cases {
		require.NoError(t, os.Chdir(c.rootModulePath), c.name)
		moduleConfigs, err := NewModuleConfigs(c.rootModulePath)
		require.NoError(t, err, c.name)
		moduleSources, err := RefactorDataSourceType(moduleConfigs, c.oldType, c.newType, c.rootModulePath)
		require.NoError(t, err, c.name)
		require.Equal(t, c.expect, moduleSources, c.name)
	}
}

func TestRefactorDataSourceName(t *testing.T) {
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	testdataDir, _ := filepath.Abs(filepath.Join(cwd, "testdata"))
	cases := []struct {
		name           string
		rootModulePath string
		dsType         string
		oldName        string
		newName        string
		expect         ModuleSources
	}{
		{
			name:           "basic",
			rootModulePath: filepath.Join(testdataDir, "data_source_name"),
			dsType:         "foo",
			oldName:        "name1",
			newName:        "name2",
			expect: map[string][]byte{
				filepath.Join(testdataDir, "data_source_name", "main.tf"): []byte(`
data "foo" "name2" {
  name = "bar"
}
`),
				filepath.Join(testdataDir, "data_source_name", "main2.tf"): []byte(`
data "foo" "namex" {
  name = data.foo.name2.name
}
`),
			},
		},
	}

	for _, c := range cases {
		require.NoError(t, os.Chdir(c.rootModulePath), c.name)
		moduleConfigs, err := NewModuleConfigs(c.rootModulePath)
		require.NoError(t, err, c.name)
		moduleSources, err := RefactorDataSourceName(moduleConfigs, c.dsType, c.oldName, c.newName, c.rootModulePath)
		require.NoError(t, err, c.name)
		require.Equal(t, c.expect, moduleSources, c.name)
	}
}

func TestRefactorDataSourceAttribute(t *testing.T) {
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	testdataDir, _ := filepath.Abs(filepath.Join(cwd, "testdata"))
	cases := []struct {
		name           string
		rootModulePath string
		dsType         string
		dsName         string
		oldAddr        []string
		newAddr        []string
		expect         ModuleSources
	}{
		{
			name:           "basic",
			rootModulePath: filepath.Join(testdataDir, "data_source_attribute"),
			dsType:         "foo",
			dsName:         "a",
			oldAddr:        []string{"addr1"},
			newAddr:        []string{"addr2"},
			expect: map[string][]byte{
				filepath.Join(testdataDir, "data_source_attribute", "main.tf"): []byte(`
data "foo" "a" {
  addr2 = "x"
}
`),
				filepath.Join(testdataDir, "data_source_attribute", "main2.tf"): []byte(`
data "foo" "b" {
  name = data.foo.a.addr2
}
`),
			},
		},
	}

	for _, c := range cases {
		require.NoError(t, os.Chdir(c.rootModulePath), c.name)
		moduleConfigs, err := NewModuleConfigs(c.rootModulePath)
		require.NoError(t, err, c.name)
		moduleSources, err := RefactorDataSourceAttribute(moduleConfigs, c.dsType, c.dsName, c.oldAddr, c.newAddr, c.rootModulePath)
		require.NoError(t, err, c.name)
		require.Equal(t, c.expect, moduleSources, c.name)
	}
}
