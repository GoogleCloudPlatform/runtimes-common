/*
Library for parsing versions.yaml file.
*/
package versions

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Package struct {
	Version      string
	Major        string
	Gpg          string
	Sha1         string
	Sha256       string
	Md5          string
	RetryCommand string
}

type Version struct {
	Dir          string
	Repo         string
	Tags         []string
	From         string
	Cmd          string
	Packages     map[string]Package
	ExcludeTests []string `yaml:"excludeTests"`
}

type Spec struct {
	Versions   []Version
	Extensions []string
	Config     map[string]string
}

func LoadVersions(path string) Spec {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	spec := Spec{}
	err = yaml.Unmarshal([]byte(data), &spec)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	validateUniqueTags(spec)
	return spec
}

func validateUniqueTags(spec Spec) {
	repoTags := make(map[string]bool)
	for _, version := range spec.Versions {
		for _, tag := range version.Tags {
			repoTag := fmt.Sprintf("%s:%s", version.Repo, tag)
			if repoTags[repoTag] {
				log.Fatalf("error: duplicate repo tag %v", repoTag)
			}
			repoTags[repoTag] = true
		}
	}
}

func (spec *Spec) CheckExtension(s string) bool {
	for _, b := range spec.Extensions {
		if b == s {
			return true
		}
	}
	return false
}
