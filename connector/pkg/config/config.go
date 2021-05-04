package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type JSONConfig struct {
	Connection struct {
		URL                 string `yaml:"url"`
		DB                  string `yaml:"db"`
		Query               string `yaml:"query"`
		Fcq                 string `yaml:"fcq"`
		DownloadListMaxsize int    `yaml:"download_list_maxsize"`
		DownloadBatchSize   int    `yaml:"download_batch_size"`
	} `yaml:"connection"`

	Output struct {
		ConvertEnable bool `yaml:"convert_enable"`

		Elasticsearch struct {
			Enable bool   `yaml:"enable"`
			Host   string `yaml:"host"`
			Login  string `yaml:"login"`
			Pwd    string `yaml:"pwd"`
			Index  string `yaml:"index"`
		} `yaml:"elasticsearch"`

		FileSystem struct {
			Enable  bool   `yaml:"enable"`
			JSONDir string `yaml:"json_dir"`
		} `yaml:"filesystem"`

		CsvFile struct {
			Enable   bool   `yaml:"enable"`
			FileName string `yaml:"file_name"`
		} `yaml:"csv_file"`
	} `yaml:"output"`
}

type PDFConfig struct {
	DB      string `yaml:"db"`
	Dir     string `yaml:"dir"`
	CsvFile string `yaml:"csv_file"`

	Auth struct {
		ASPXAUTH        string `yaml:".ASPXAUTH"`
		ASPNETSessionID string `yaml:"ASP.NET_SessionId"`
	} `yaml:"auth"`
}

type Config struct {
	JSONConfig JSONConfig `yaml:"json"`
	PDFConfig  PDFConfig  `yaml:"pdf"`
}

func NewConfig(configPath string) (*Config, error) {
	s, err := os.Stat(configPath)
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		return nil, fmt.Errorf("'%s' is a directory, not a normal file", configPath)
	}

	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
