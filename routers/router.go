package routers

import (
	_ "github.com/betterDuanjiawei/gin-jianyu/docs"
	"github.com/betterDuanjiawei/gin-jianyu/middleware/jwt"
	"github.com/betterDuanjiawei/gin-jianyu/pkg/export"
	"github.com/betterDuanjiawei/gin-jianyu/pkg/setting"
	"github.com/betterDuanjiawei/gin-jianyu/pkg/upload"
	"github.com/betterDuanjiawei/gin-jianyu/routers/api"
	v1 "github.com/betterDuanjiawei/gin-jianyu/routers/api/v1"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))

	r.GET("/auth", api.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload", api.UploadImage)

	//导出标签
	r.POST("/tags/export", v1.ExportTag)
	//导入标签
	r.POST("/tags/import", v1.ImportTag)

	apiv1 := r.Group("api/v1")
	apiv1.Use(jwt.JWT())
	{
		apiv1.GET("/tags", v1.GetTags)
		apiv1.POST("/tags", v1.AddTag)
		apiv1.PUT("tags/:id", v1.EditTag)
		apiv1.DELETE("tags/:id", v1.DeleteTag)

		apiv1.GET("/articles", v1.GetArticles)
		apiv1.GET("/articles/:id", v1.GetArticle)
		apiv1.POST("/articles", v1.AddArticle)
		apiv1.PUT("/articles/:id", v1.EditArticle)
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
