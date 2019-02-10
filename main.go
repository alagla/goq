package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const fname1 = "C:/Users/evaldas/Documents/proj/Java/github.com/qupla/src/main/resources/Qupla.yml"

//const fname2 = "C:/Users/evaldas/Documents/proj/site_data/test.yml"
const fname = fname1

const testout = "C:/Users/evaldas/Documents/proj/site_data/tmp/echotest.yml"

func main() {
	var quplaYAML QuplaModuleYAML
	must(readYAML(fname, &quplaYAML))
	fmt.Printf("%+v", quplaYAML)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func readYAML(fname string, outStruct interface{}) error {
	fmt.Printf("reading %v\n", fname)
	yamlFile, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer yamlFile.Close()

	yamlbytes, _ := ioutil.ReadAll(yamlFile)

	err = yaml.Unmarshal(yamlbytes, outStruct)
	if err != nil {
		return err
	}
	outData, err := yaml.Marshal(&outStruct)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(testout, outData, 0644)
	return nil
}
