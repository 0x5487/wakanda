package config

import (
	"flag"
	"io/ioutil"
	"log"
	stdlog "log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type LogSetting struct {
	Name             string `yaml:"name"`
	Type             string `yaml:"type"`
	MinLevel         string `yaml:"min_level"`
	ConnectionString string `yaml:"connection_string"`
}

type Configuration struct {
	Env      string
	Logs     []LogSetting `yaml:"logs"`
	Database struct {
		Username string
		Password string
		Address  string
		Type     string
		DBName   string
	}
	Messenger struct {
		Bind string
	}
	Gateway struct {
		Bind string
	}
}

func New(fileName string) *Configuration {
	flag.Parse()

	c := Configuration{}

	//read and parse config file
	rootDirPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		stdlog.Fatalf("config: file error: %s", err.Error())
	}

	configPath := filepath.Join(rootDirPath, fileName)
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		// config exists
		file, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Fatalf("config: file error: %s", err.Error())
		}

		err = yaml.Unmarshal(file, &c)
		if err != nil {
			log.Fatal("config: config error:", err)
		}
	}

	return &c
}
