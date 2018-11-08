package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

var once sync.Once
var yconfig *Config

type MachineDetails struct {
	IP   string `yaml:"ip"`
	PORT string `yaml:"port"`
}

type MySQLDetails struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	PoolSize int    `yaml:"poolsize"`
}

type Config struct {
	Server  MachineDetails `yaml:"server"`
	Client  MachineDetails `yaml:"client"`
	DB      MySQLDetails   `yaml:"db"`
	LogFile string         `yaml:"logfile"`
}

func getYamlConfig(filename string) *Config {
	once.Do(func() {
		yamlFile, err := ioutil.ReadFile(filename)

		if err != nil {
			panic(err)
		}
		var configd Config
		err = yaml.Unmarshal(yamlFile, &configd)

		if err != nil {
			panic(err)
		}
		yconfig = &configd
	})
	return yconfig
}
