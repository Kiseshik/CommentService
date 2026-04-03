package port

import (
	"github.com/gin-gonic/gin"
)

type Handler interface {
	RegisterPublicHandlers(group *gin.RouterGroup)
	RegisterPrivateHandlers(group *gin.RouterGroup)
	RegisterInternalHandlers(group *gin.RouterGroup)
}
