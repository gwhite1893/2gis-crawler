package config

import (
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTPServer *HTTPServerCfg `yaml:"http_server"`
	Logger     *logrus.Logger `yaml:"log"`
	Tracer     *TracingConfig `yaml:"tracing"`
	Crawler    *CrawlerCfg    `yaml:"crawler"`
}

type CrawlerCfg struct {
	RequestTimeOutSec int `yaml:"request_time_out_sec"`
}

type HTTPServerCfg struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (s *HTTPServerCfg) Addr() string {
	return fmt.Sprintf("%v:%v", s.Host, s.Port)
}

type TracingConfig struct {
	URL            string  `yaml:"url"`
	SampleFraction float64 `yaml:"sample_fraction"`
}

type URL struct {
	*url.URL
}

func (u *URL) UnmarshalYAML(unmarshal func(interface{}) error) error {
	stringURL := ""

	err := unmarshal(&stringURL)
	if err != nil {
		return err
	}

	u.URL, err = url.Parse(stringURL)

	return err
}

func New(filePath string) (*Config, error) {
	c := Config{}

	configContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(configContent, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
