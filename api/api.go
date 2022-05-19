package api

import (
	"time"

	_ "github.com/merico-dev/lake/api/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/config"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title  DevLake Swagger API
// @version 0.1
// @description  <h2>This is the main page of devlake api</h2>
// sdfasdfasd
// @license.name Apache-2.0
// @license.url https://www.baidu.com
// @host localhost:8080
// @BasePath /
func CreateApiService() {
	v := config.GetConfig()
	gin.SetMode(v.GetString("MODE"))
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// CORS CONFIG
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           120 * time.Hour,
	}))

	RegisterRouter(router)
	err := router.Run(v.GetString("PORT"))
	if err != nil {
		panic(err)
	}
}
