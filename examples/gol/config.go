package main

import (
	. "github.com/lunfardo314/goq/cfg"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type golConfigStruct struct {
	WebServerPort   int    `yaml:"webServerPort"`
	SiteRoot        string `yaml:"siteRoot"`
	GolCodeLocation string `yaml:"golCodeLocation"`
}

var ConfigGol = golConfigStruct{}

func readConfig(fname string) error {
	Logf(0, "reading config: %v", fname)
	yamlFile, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer yamlFile.Close()

	yamlbytes, _ := ioutil.ReadAll(yamlFile)

	err = yaml.Unmarshal(yamlbytes, &ConfigGol)
	if err != nil {
		return err
	}
	return nil
}
