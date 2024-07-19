package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) getAllMessage(c *gin.Context) {
	personId, err := getPersonId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	roomIdStr, ok := c.GetQuery("room_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "should be room_id"})
		return
	}
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for room_id"})
		return
	}
	dateStr, ok := c.GetQuery("date")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "should be date"})
		return
	}
	date, err := parseDate(dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}
	limitStr := c.Query("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit format"})
		return
	}
	messages, err := h.service.MessageService.GetAllMessageForChatRoom(personId, roomId, date, limit)
	if err != nil {
		log.Printf("error loading messages for chat room: %s", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": "server error"})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func parseDate(date string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
