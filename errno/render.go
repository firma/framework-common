package errno

import (
	"github.com/gin-gonic/gin"
)

type (
	Render interface {
		// DEPRECATED
		RenderGinJson(ctx *gin.Context)
		RenderJson(ctx *gin.Context)

		// 此方法只针对自定义Data直接输出Json
		RenderCustomJson(ctx *gin.Context)

		// 此方法只针对string/[]byte类型的data直接输出text/plain
		RenderText(ctx *gin.Context)
	}
)

func (e errno) RenderJson(ctx *gin.Context) {
	ctx.JSON(e.GetHttpStatusCode(), e.reset(ctx))
	ctx.Abort()
}

func (e errno) RenderCustomJson(ctx *gin.Context) {
	ctx.JSON(e.GetHttpStatusCode(), e.GetData())
	ctx.Abort()
}

func (e errno) RenderText(ctx *gin.Context) {
	ctx.Data(e.GetHttpStatusCode(), "text/plain", e.GetRawData())
	ctx.Abort()
}

// DEPRECATED
func (e errno) RenderGinJson(ctx *gin.Context) {
	e.RenderJson(ctx)
}
