package router

import (
	"fmt"
	"os"

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
	// 초기 설치시 환경변수 설정 필요(dev 또는 release)
	// os.Setenv("WBABEProjectMode","dev")
	serverMode := os.Getenv("WBABEProjectMode")
	if serverMode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else if serverMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		panic("WBABEProjectMode 환경변수 설정이 필요합니다!(router 참조)")
	}
	fmt.Printf("server mode: %s \n", serverMode)
	e := gin.New()
	e.Use(gin.Logger())
	e.Use(gin.Recovery())
	e.Use(logger.GinLogger())
	e.Use(logger.GinRecovery(true))
	e.Use(CORS())
	//swagger 핸들러 미들웨어에 등록
	e.GET("/swagger/:any", ginSwg.WrapHandler(swgFiles.Handler))
	docs.SwaggerInfo.Host = "localhost:8080"

	e.POST("/user", p.ct.InsertUserControl) // 유저 등록 주소로 처음 환경 초기화시 사용을 위해 추가
	e.POST("/test", p.ct.TestControl)       // 일부 작은 단위의 기능을 테스트하기 위한 주소로 swagger에 등록하지 않음

	owner := e.Group("owner", liteAuth())
	{
		owner.GET("/menu", p.ct.GetMenuControl) // 임시로 GetOk로 연결
		owner.GET("/menu/detail", p.ct.GetMenuDetailControl)
		owner.POST("/menu", p.ct.InsertMenuControl)
		owner.PUT("/menu", p.ct.UpdateMenuControl)
		owner.DELETE("/menu", p.ct.DeleteMenuControl)

		owner.GET("/order", p.ct.GetOrderControl)         // 오더 상태 확인
		owner.PUT("/order", p.ct.UpdateOwnerOrderControl) // 오더 상태 수정
	}

	customer := e.Group("customer", liteAuth())
	{
		customer.GET("/menu", p.ct.GetMenuControl)
		customer.GET("/menu/detail", p.ct.GetMenuDetailControl) // 개별 메뉴에 대한 평점 및 리뷰 데이터 확인

		customer.GET("/order", p.ct.GetOrderControl)             // 자신이 주문한 order 확인
		customer.POST("/order", p.ct.InsertCustomerOrderControl) // order 주문
		customer.PUT("/order", p.ct.UpdateCustomerOrderControl)  // 자신이 주문한 order 정보 변경, 삭제 대신 취소 상태로 대신함

		customer.GET("/order/review", p.ct.GetReviewControl)       // 자신이 주문한 order에 대한 리뷰 확인
		customer.POST("/order/review", p.ct.InsertReviewControl)   // 자신이 주문한 order에 대한 리뷰 추가
		customer.PUT("/order/review", p.ct.UpdateReviewControl)    // 자신이 주문한 order에 대한 리뷰 수정
		customer.DELETE("/order/review", p.ct.DeleteReviewControl) // 자신이 주문한 order에 대한 리뷰 삭제
	}

	return e
}
