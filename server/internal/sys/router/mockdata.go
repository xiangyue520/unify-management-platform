package router

import (
	"mayfly-go/internal/sys/api"
	"mayfly-go/pkg/ctx"

	"github.com/gin-gonic/gin"
)

func InitMockRouter(router *gin.RouterGroup) {
	r := &api.MockData{}
	db := router.Group("mock-datas")
	{

		db.GET("", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(r.MockDatas)
		})

		db.GET("/services", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(r.MockDataServices)
		})

		db.GET("/detail", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(r.MockDataDetail)
		})

		db.POST("", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).
				Handle(r.SaveMockData)
		})

		db.DELETE(":id", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(r.DelMockData)
		})

	}
}
