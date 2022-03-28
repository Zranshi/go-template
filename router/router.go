package router

import (
	"go-template/config"

	"github.com/gin-gonic/gin"
)

var Router = gin.Default()

func init() {
	if config.Conf.GetString("host.mode") == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}
