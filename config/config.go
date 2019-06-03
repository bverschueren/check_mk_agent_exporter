package config

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Target struct {
	HostName     string `yaml:"HostName"`
	Port         int    `yaml:"Port"`
	User         string `yaml:"User"`
	IdentityFile string `yaml:"IdentityFile"`
}

type Config struct {
	Filename *string
}

func (c Config) ReadFile(targets *map[string]Target) {
	source, err := ioutil.ReadFile(*c.Filename)
	if err != nil {
		log.Fatalf("Unable to open '%s': %s", *c.Filename, err)
	}
	// read from 'targets' root element
	targetlist := struct {
		List *map[string]Target `yaml:"targets"`
	}{
		targets,
	}
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &targetlist)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Debugf("targets: %+v", targetlist.List)
}

func (t *Target) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawTarget Target
	raw := rawTarget{
		Port:         22,
		IdentityFile: "~/.ssh/id_rsa",
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*t = Target(raw)
	return nil
}
