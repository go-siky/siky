package config

import (
	"errors"
	"os"

	"github.com/qiniu/log"
	"gopkg.in/yaml.v2"
)

//Server config
type Server struct {
	Port string
}

//Database config
type Database struct {
	Mongo
}

//Mongo database config
type Mongo struct {
	URL string
}

//Config for project
type Config struct {
	Server
	Database
	Logger struct {
		Level int
	}
	Registry struct {
		Protocol string
		URL      string
		Version  string
	}
}

var cfg *Config

//New for project
func New(basePath string) (*Config, error) {

	cfg = &Config{}
	log.Println(basePath + "config.yml")
	//open config file
	cfgFile, err := os.Open(basePath + "config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer cfgFile.Close()

	cfgBufLen, err := cfgFile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	cfgFBuf := make([]byte, cfgBufLen.Size())

	cfgFile.Read(cfgFBuf)
	if len(cfgFBuf) == 0 {
		return nil, errors.New("not config")
	}
	yaml.Unmarshal(cfgFBuf, &cfg)
	log.Info("%v", cfg)
	return cfg, nil
}

// Get config
func Get() *Config {
	if cfg == nil {
		log.Error("Config not initialized")
	}
	return cfg
}
