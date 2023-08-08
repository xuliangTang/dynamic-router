package lib

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

type Config struct {
	Routes RouteConfigs `json:"routes"`
}

// SaveToFile 重新写入文件
func (this *Config) SaveToFile() {
	b, err := yaml.Marshal(this)
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.OpenFile("./config.yaml", os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		log.Fatalln("2", err)
	}
}

type RouteConfig struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Code   int    `json:"code"`
}

// Register 注册路由
// append 是否追加写入到配置文件
func (this *RouteConfig) Register(isAppend bool) {
	if !this.Exist() {
		GinServer.Handle(this.Method, this.Path, func(ctx *gin.Context) {
			ctx.JSON(this.Code, gin.H{"msg:": fmt.Sprintf("%s %s", this.Method, this.Path)})
		})

		if isAppend {
			SysConfig.Routes = append(SysConfig.Routes, this)
			SysConfig.SaveToFile()
		}
	}
}

// Exist 判断路由是否存在
func (this *RouteConfig) Exist() bool {
	return hasRoute(this.Method, this.Path, GinServer.Routes())
}

type RouteConfigs []*RouteConfig

// 读取配置文件
func loadConfigs() *Config {
	f, err := os.Open("./config.yaml")
	if err != nil {
		log.Fatalln("读取配置文件出错:", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatalln("读取配置文件内容出错:", err)
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(b, cfg); err != nil {
		log.Fatalln("配置文件反序列化出错:", err)
	}

	return cfg
}

// 把配置文件的routes注册到路由
func registerRoutes() {
	if SysConfig.Routes != nil {
		for _, r := range SysConfig.Routes {
			r.Register(false)
		}
	}
}
