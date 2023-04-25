package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"template/pkg/frame"
	"template/pkg/middleware"
)

func main() {
	defer glog.Flush()

	frame.InitFramework()
	defer frame.UnInitFramework()

	debug := viper.GetBool("application.debug")
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	//r.Use(static.Serve("/", static.LocalFile("./bin/dist", true)))

	r.Use(middleware.Cors())
	r.Use(middleware.LoggerToFile())
}
