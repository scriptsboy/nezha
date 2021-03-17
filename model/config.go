package model

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	ConfigTypeGitHub = "github"
	ConfigTypeGitee  = "gitee"
)

type Config struct {
	Debug bool
	Site  struct {
		Brand        string // 站点名称
		CookieName   string // 浏览器 Cookie 名称
		Theme        string
		CustomCode   string
		ViewPassword string // 前台查看密码
	}
	Oauth2 struct {
		Type         string
		Admin        string // 管理员用户名列表
		ClientID     string
		ClientSecret string
	}
	HTTPPort                   uint
	GRPCPort                   uint
	EnableIPChangeNotification bool

	v *viper.Viper
}

func (c *Config) Read(path string) error {
	c.v = viper.New()
	c.v.SetConfigFile(path)
	err := c.v.ReadInConfig()
	if err != nil {
		return err
	}

	err = c.v.Unmarshal(c)
	if err != nil {
		return err
	}

	if c.Site.Theme == "" {
		c.Site.Theme = "default"
	}

	c.v.OnConfigChange(func(in fsnotify.Event) {
		c.v.Unmarshal(c)
		fmt.Println("配置文件更新，重载配置", c)
	})

	go c.v.WatchConfig()
	return nil
}

func (c *Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.v.ConfigFileUsed(), data, os.ModePerm)
}
