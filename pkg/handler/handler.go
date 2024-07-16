package handler

import (
	"chatService/pkg/handler/ws"
	"chatService/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	hub := ws.NewHub()
	go hub.Run()
	router.GET("/ws", h.personIdentity, h.handleWs(hub))
	auth := router.Group("/auth")
	{
		auth.POST("/registration", h.createPerson)
		auth.POST("/login", h.logIn)
	}
	api := router.Group("/api", h.personIdentity)
	{
		room := api.Group("/room")
		{
			room.GET("", h.getAllRoom)
			room.POST("", h.addChatRoom)
			room.DELETE("", h.deleteChatRoom)
			room.DELETE("/exit", h.exitFromChatRoom)
			room.DELETE("/person", h.deletePersonFromChatRoom)
		}
	}
	return router
}
