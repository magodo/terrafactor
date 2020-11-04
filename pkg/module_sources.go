package pkg

// Referred from terraform/configs/configupgrade/module_sources.go

import (
	"github.com/hashicorp/terraform/configs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type ModuleSources map[string][]byte

// LoadModule looks for Terraform configuration files in the given directory
// and loads each of them into memory as source code, in preparation for
// further analysis and conversion.
//
// At this stage the files are not parsed at all. Instead, we just read the
// raw bytes from the file so that they can be passed into a parser in a
// separate step.
//
// If the given directory or any of the files cannot be read, an error is
// returned. It is not safe to proceed with processing in that case because
// we cannot "see" all of the source code for the configuration.
func LoadModule(dir string) (ModuleSources, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	ret := make(ModuleSources)
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			continue
		}
		if configs.IsIgnoredFile(name) {
			continue
		}
		if !strings.HasSuffix(name, ".tf") {
			continue
		}

		fullPath := filepath.Join(dir, name)
		src, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return nil, err
		}

		ret[fullPath] = src
	}

	return ret, nil
}
