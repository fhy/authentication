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
	Server   config.Server       `yaml:"server"`
	Log      config.LogConfig    `yaml:"log"`
	Redis    config.RedisConfig  `yaml:"redis"`
	WeChat   WechatConfig        `yaml:"wechat"`
	Cookie   config.CookieConfig `yaml:"cookie"`
	DbType   string              `yaml:"dbtype"`
	Dbconfig config.SqliteConfig `yaml:"db"`
	Jwt      config.Jwt          `yaml:"jwt"`
}

var (
	Conf *Config
)

func Init(env string) {
	configFile := "./conf/setting_local.yml"
	authConfigFile := "./conf/auth.yml"
	switch env {
	case "develop":
		configFile = "./conf/setting_dev.yml"
	case "release":
		configFile = "./conf/setting.yml"
	default:
		configFile = "./conf/setting_local.yml"
	}
	var err error
	cfg := []byte{}
	configFiles := []string{authConfigFile, configFile}
	for _, configFile := range configFiles {
		if config := LoadConf(configFile); err != nil {
			fmt.Print("Decode Config Error", err)
			os.Exit(1)
		} else {
			cfg = append(cfg, *config...)
		}
		cfg = append(cfg, []byte("\n")...)
	}
	if Conf, err = LoadConfYaml(cfg); err != nil {
		fmt.Print("Decode Config Error", err)
		os.Exit(1)
	}
	fmt.Println(Conf)
}

func LoadConf(configPath string) *[]byte {
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
	return &configFile
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
