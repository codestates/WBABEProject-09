package router

import (
	"fmt"

	ctl "WBABEProject-09/controller"
	"WBABEProject-09/docs"
	"WBABEProject-09/logger"

	"github.com/gin-gonic/gin"
	swgFiles "github.com/swaggo/files"
	ginSwg "github.com/swaggo/gin-swagger"
)

type Router struct {
	ct *ctl.Controller
}

func NewRouter(ctl *ctl.Controller) (*Router, error) {
	r := &Router{ct: ctl} //controller 포인터를 ct로 복사, 할당

	return r, nil
}

// cross domain을 위해 사용
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Forwarded-For, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// 임의 인증을 위한 함수
func liteAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c == nil {
			c.Abort()
			return
		}
		auth := c.GetHeader("Authorization")
		fmt.Println("Authorization-word", auth)
		c.Next()
	}
}

// 실제 라우팅
func (p *Router) Idx() *gin.Engine {
	e := gin.New()

	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.Use(logger.GinLogger())
	e.Use(logger.GinRecovery(true))
	e.Use(CORS())
	//swagger 핸들러 미들웨어에 등록
	e.GET("/swagger/:any", ginSwg.WrapHandler(swgFiles.Handler))
	docs.SwaggerInfo.Host = "localhost"

	owner := e.Group("owner", liteAuth())
	{
		owner.GET("/menu", p.ct.GetOK) // 임시로 GetOk로 연결
		owner.POST("/menu", p.ct.InsertMenuControl)
		owner.PUT("/menu", p.ct.GetOK)
		owner.DELETE("/menu", p.ct.GetOK)
	}

	customer := e.Group("customer", liteAuth())
	{
		customer.GET("/menu", p.ct.GetOK)
	}

	return e
}