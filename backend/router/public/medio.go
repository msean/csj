package public

import (
	"github.com/gin-gonic/gin"
)

type PublicRouter struct{}

func (s *PublicRouter) InitMedioRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {

	routerWithoutAuth := PublicRouter.Group("public")
	{
		routerWithoutAuth.POST("/uploadMedia", medioApi.UploadMedia)
		routerWithoutAuth.GET("/download", medioApi.Download)
	}
}
