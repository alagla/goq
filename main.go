package main

import (
	"github.com/lunfardo314/goq/program"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const fname = "C:/Users/evaldas/Documents/proj/Java/github.com/qupla/src/main/resources/Qupla.yml"
const testout = "C:/Users/evaldas/Documents/proj/site_data/tmp/echotest.yml"

func main() {
	quplaModule := program.NewQuplaModule()
	must(readYAML(fname, &quplaModule))
	if !quplaModule.Analyze() {
		errorf("Failed analyzing Qupla module")
	}
	quplaModule.PrintStats()
	if err := quplaModule.Execs[0].Execute(); err != nil {
		errorf("error: %v", err)
	}
	infof("Ciao!")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func readYAML(fname string, outStruct interface{}) error {
	infof("reading %v", fname)
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
