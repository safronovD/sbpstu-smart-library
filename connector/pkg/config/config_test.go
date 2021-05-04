package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	yaml_file = "/tmp/config.yaml"
)

func TestNewConfig(t *testing.T) {
	// Create Config struct and set some fields
	config := &Config{}

	config.JSONConfig.Connection.DB = "some_data"
	config.JSONConfig.Output.FileSystem.JSONDir = "some_data2"
	config.PDFConfig.Auth.ASPNETSessionID = "some_data3"

	// Write data to file
	data, err := yaml.Marshal(config)
	assert.Nil(t, err)

	err = ioutil.WriteFile(yaml_file, data, 0644)
	assert.Nil(t, err)

	// Read data from file with NewConfig method and check structs are the same
	readConnfig, err := NewConfig(yaml_file)
	assert.Nil(t, err)
	assert.Equal(t, config, readConnfig)
}
