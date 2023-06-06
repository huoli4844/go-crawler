package config

import (
	"github.com/lizongying/go-crawler/pkg/cli"
	"gopkg.in/yaml.v2"
	"log"
	"net/url"
	"os"
	"time"
)

const defaultHttpProto = "2.0"
const defaultTimeout = time.Minute
const defaultDevServer = "http://:8081"

type Config struct {
	MongoEnable bool `yaml:"mongo_enable" json:"-"`
	Mongo       struct {
		Example struct {
			Uri      string `yaml:"uri" json:"-"`
			Database string `yaml:"database" json:"-"`
		} `yaml:"example" json:"-"`
	} `yaml:"mongo" json:"-"`
	MysqlEnable bool `yaml:"mysql_enable" json:"-"`
	Mysql       struct {
		Example struct {
			Uri      string `yaml:"uri" json:"-"`
			Database string `yaml:"database" json:"-"`
		} `yaml:"example" json:"-"`
	} `yaml:"mysql" json:"-"`
	KafkaEnable bool `yaml:"kafka_enable" json:"-"`
	Kafka       struct {
		Example struct {
			Uri string `yaml:"uri" json:"-"`
		} `yaml:"example" json:"-"`
	} `yaml:"kafka" json:"-"`
	Log struct {
		Filename string `yaml:"filename" json:"-"`
		LongFile bool   `yaml:"long_file" json:"-"`
		Level    string `yaml:"level" json:"-"`
	} `yaml:"log" json:"-"`
	Proxy struct {
		Example string `yaml:"example" json:"-"`
	} `yaml:"proxy" json:"-"`
	Request struct {
		Concurrency   int    `yaml:"concurrency" json:"-"`
		Interval      int    `yaml:"interval" json:"-"`
		Timeout       int    `yaml:"timeout" json:"-"`
		OkHttpCodes   []int  `yaml:"ok_http_codes" json:"-"`
		RetryMaxTimes int    `yaml:"retry_max_times" json:"-"`
		HttpProto     string `yaml:"http_proto" json:"-"`
	} `yaml:"request" json:"-"`
	DevServer string `yaml:"dev_server" json:"-"`
}

func (c *Config) GetProxy() *url.URL {
	if c.Proxy.Example != "" {
		proxy, err := url.Parse(c.Proxy.Example)
		if err != nil {
			log.Panicln(err)
		}
		return proxy
	}

	return nil
}

func (c *Config) GetHttpProto() string {
	if c.Request.HttpProto != "" {
		return c.Request.HttpProto
	}

	return defaultHttpProto
}

func (c *Config) GetTimeout() time.Duration {
	if c.Request.Timeout > 0 {
		return time.Second * time.Duration(c.Request.Timeout)
	}

	return defaultTimeout
}

func (c *Config) GetDevServer() (url *url.URL, err error) {
	if c.DevServer != "" {
		url, err = url.Parse(c.DevServer)
		return
	}

	url, err = url.Parse(defaultDevServer)
	return
}

func (c *Config) LoadConfig(configPath string) (err error) {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Panicln(err)
	}

	err = yaml.Unmarshal(configData, c)
	if err != nil {
		log.Panicln(err)
	}

	return
}

func NewConfig(cli *cli.Cli) (config *Config, err error) {
	config = &Config{}
	configFile := cli.ConfigFile
	if configFile != "" {
		err = config.LoadConfig(configFile)
		if err != nil {
			log.Panicln(err)
		}
	}

	return
}
