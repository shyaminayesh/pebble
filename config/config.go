package config

import (
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v3"
)


type (
	Configuration struct {
		Connection	Connection
		Schema		Schema
		Backup		Backup
	}

	Connection struct {
		Host		string		`yaml:"host"`
		Port		uint		`yaml:"port"`
		Database	string		`yaml:"database"`
		User		string		`yaml:"user"`
		Password	string		`yaml:"password"`
		Dialect		string		`yaml:"dialect"`
	}

	Schema struct {
		Directory	string		`yaml:"directory"`
	}

	Backup struct {
		File		File		`yaml:"file"`
		Directory	string		`yaml:"directory"`
	}

	File struct {
		Timestamp	string		`yaml:"timestamp"`
	}
)


func Get() Configuration {

	/**
	* We need to read the file and append the configuration
	* properties to the Configuration struct to continue
	* the configuration file parse
	*/
	buffer, err := ioutil.ReadFile("./pebble.yml")
	if err != nil { log.Fatal(err) }

	/**
	* We have to file ready to parse using the yaml
	* parser. We can return the Configuration struct
	* once we complete parse.
	*/
	configuration := Configuration {}
	yaml.Unmarshal(buffer, &configuration)
	return configuration

}
