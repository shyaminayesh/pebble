package config

import (
	"fmt"
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v3"
	"github.com/spf13/viper"
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
		Directory	string		`yaml:"directory"`
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
	fmt.Println( configuration )
	return configuration

}


func Config() *viper.Viper {

	var v = viper.New()
	v.SetConfigName("pebble")
	v.SetConfigType("yml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return v
}