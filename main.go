package main

import (
	"github.com/gin-gonic/gin"

	"github.com/gamezop/util"
)

func main() {
	service := gin.Default()
	initService(service)
	service.Run("localhost:5000")
}

//initService : API Service
func initService(group *gin.Engine) {
	service := group.Group("api")
	{
		service.POST("/", PostData)
		service.GET("/", GetData)
	}
}

//PostData : Post API
func PostData(c *gin.Context) {
	var err error
	if err = util.PushDataToRedis(c); err == nil {
		c.JSON(200, "success")
	} else {
		c.JSON(200, err)
	}
}

//GetData : Get API
func GetData(c *gin.Context) {
	var (
		data interface{}
		err  error
	)
	if data, err = util.GetDataFromDatabase(c); err == nil {
		c.JSON(200, data)
	} else {
		c.JSON(200, err)
	}
}
