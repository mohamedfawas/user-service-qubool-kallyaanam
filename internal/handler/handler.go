package handler

import "github.com/gin-gonic/gin"

// Handler defines the interface for all API handlers
type Handler interface {
	RegisterRoutes(router *gin.RouterGroup)
}
