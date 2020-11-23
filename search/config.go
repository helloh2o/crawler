package search

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var (
	Config *conf
)

type conf struct {
	EsNode     string              `yaml:"es_node"`
	Addr       string              `yaml:"addr"`
	IndexQuery map[string][]string `yaml:"index_query"`
	IndexKey   map[string]string   `yaml:"-"`
}

func Init(filename string) {
	Config = &conf{IndexKey: make(map[string]string)}
	if yamlFile, err := ioutil.ReadFile(filename); err != nil {
		log.Fatal(err)
	} else if err = yaml.Unmarshal(yamlFile, Config); err != nil {
		log.Fatal(err)
	}
	for index, _ := range Config.IndexQuery {
		Config.IndexKey[index] = index
	}
	log.Printf("config::\n %+v", Config)
}
