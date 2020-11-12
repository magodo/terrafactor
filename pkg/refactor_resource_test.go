package pkg

import (
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
			name:           "basic",
			rootModulePath: filepath.Join(testdataDir, "resource_attribute"),
			resType:        "foo",
			resName:        "a",
			oldAddr:        []string{"addr1"},
			newAddr:        []string{"addr2"},
			expect: map[string][]byte{
				filepath.Join(testdataDir, "resource_attribute", "main.tf"): []byte(`
resource "foo" "a" {
  addr2 = "x"
}
`),
				filepath.Join(testdataDir, "resource_attribute", "main2.tf"): []byte(`
resource "foo" "b" {
  name = foo.a.addr2
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
		require.Equal(t, c.expect, moduleSources, c.name)
	}
}
