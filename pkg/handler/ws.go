package handler

import (
	"chatService/pkg/handler/ws"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) handleWs(hub *ws.Hub) func(c *gin.Context) {
	return func(c *gin.Context) {
		personId, err := getPersonId(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		roomIdStr, ok := c.GetQuery("room_id")
		if !ok {
			c.JSON(http.StatusBadRequest, "should be query room_id")
			return
		}
		roomId, err := strconv.Atoi(roomIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		room, err := h.service.ChatRoomService.GetByPersonIdAndChatRoomId(personId, roomId)
		if room.Id == 0 && room.Name == "" {
			c.JSON(http.StatusBadRequest, "absent chat room for this person")
			return
		}
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "server error"})
			return
		}
		ws.ServeWs(hub, c.Writer, c.Request, personId, roomId, h.service)
	}
}
