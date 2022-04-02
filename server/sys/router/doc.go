package router

import (
	"mayfly-go/base/ctx"
	"mayfly-go/server/sys/api"

	"github.com/gin-gonic/gin"
)

func InitDocRouter(router *gin.RouterGroup) {
	r := &api.Doc{}
	db := router.Group("/docs")
	{

		db.GET("", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(r.Docs)
		})

		db.GET("/categories", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(r.DocCategory)
		})

		db.GET(":id", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(r.DocDetail)
		})

		db.POST("", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).
				Handle(r.SaveDoc)
		})

		db.POST("/upload", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).
				Handle(r.UploadDocImage)
		})

		db.DELETE(":id", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(r.DelDoc)
		})
	}
}
