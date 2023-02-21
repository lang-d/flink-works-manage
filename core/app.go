package core

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/land-d/flink-works-amage/router"
	"log"
	"os"
)

const (
	DEV_MODEL  string = "dev"
	PROD_MODEL string = "prod"
)

type App struct {
	Host  string
	Port  int
	Model string
	engin *gin.Engine
}

// usage 错误信息的处理
func usage() {
	_, err := fmt.Fprintf(os.Stderr, "run Options:\n")
	if err != nil {
		log.Fatalf("%s", err)
	}
	flag.PrintDefaults()
}

func NewApp() *App {
	return &App{}
}

// Parse 提交运行参数
func (this *App) Parse() {
	help := flag.Bool("h", false, "see command help")
	host := flag.String("i", "0.0.0.0", "accept visit host")
	port := flag.Int("p", 8080, "server port")
	model := flag.String("m", "", "server run mode,if you need run by dev,input dev")

	flag.Parse()
	flag.Usage = usage

	if *help {
		flag.Usage()
		os.Exit(1)
	}
	this.Host = *host
	this.Port = *port
	this.Model = *model

}

func (this *App) SetHost(host string) *App {
	this.Host = host
	return this
}

func (this *App) SetPort(port int) *App {
	this.Port = port
	return this
}

func (this *App) SetModel(model string) *App {
	this.Model = model
	return this
}

func (this *App) Run() {
	// 初始化配置

	// 检查

	// 加载引擎
	// 设置运行模式
	if this.Model == DEV_MODEL {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)
	}
	this.engin = gin.New()

	router.SetupRoute(this.engin)

	err := this.engin.Run(fmt.Sprintf("%s:%d", this.Host, this.Port))
	if err != nil {
		log.Fatal(err)
	}
}
