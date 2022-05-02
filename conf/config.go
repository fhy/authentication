package conf

import (
	"base/config"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type WechatConfig struct {
	OfficialAccount config.OfficialAccountConfig `yaml:"officalAccount"`
	MiniProgram     config.MiniProgramConfig     `yaml:"miniProGram"`
}

type Config struct {
	Server config.Server       `yaml:"server"`
	Log    config.LogConfig    `yaml:"log"`
	Redis  config.RedisConfig  `yaml:"redis"`
	WeChat WechatConfig        `yaml:"wechat"`
	Cookie config.CookieConfig `yaml:"cookie"`
	Db     config.DbConfig     `yaml:"db"`
	Jwt    config.Jwt          `yaml:"jwt"`
}

var (
	Conf *Config
)

func Init(env string) {
	configFile := "./conf/setting_dev.yml"
	switch env {
	case "local":
		configFile = "./conf/setting_local.yml"
	case "release":
		configFile = "./conf/setting.yml"
	default:
		configFile = "./conf/setting_dev.yml"
	}
	var err error
	if Conf, err = LoadConf(configFile); err != nil {
		fmt.Print("Decode Config Error", err)
		os.Exit(1)
	}
	fmt.Println(Conf)
}

func LoadConf(configPath string) (*Config, error) {
	fmt.Printf("loading config file: %s\n", configPath)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Print("can't find config path\n")
		os.Exit(1)
	} else {
		if err != nil {
			fmt.Print("Decode Config Error", err)
			os.Exit(1)
		}
	}
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Print("can't read config file\n")
		os.Exit(1)
	}
	return LoadConfYaml(configFile)
}

func LoadConfYaml(configFile []byte) (*Config, error) {
	var conf *Config
	err := yaml.Unmarshal(configFile, &conf)
	if err != nil {
		fmt.Print("can't parse config file\n")
		return nil, err
	}

	fmt.Println(conf)
	config.FormatConfig(conf)
	return conf, nil
}
