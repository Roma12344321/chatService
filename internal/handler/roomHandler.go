package handler

import (
	"chatService/pkg/model"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) getAllRoom(c *gin.Context) {
	personId, err := getPersonId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rooms, err := h.service.ChatRoomService.GetAllChatRoom(personId)
	if err != nil {
		log.Printf("error getting chat rooms: %s", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": "server error"})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

func (h *Handler) addChatRoom(c *gin.Context) {
	personId, err := getPersonId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	otherIdStr, ok := c.GetQuery("person_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "there must be other person id"})
		return
	}
	otherId, err := strconv.Atoi(otherIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "there must be other person id"})
		return
	}
	err = h.service.ChatRoomService.CreateChatRoom(personId, otherId)
	if err != nil {
		log.Printf("error creating chat rooms: %s", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": "server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "room was created"})
}

func (h *Handler) deleteChatRoom(c *gin.Context) {
	roomIdStr, ok := c.GetQuery("room_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "should be query room_id"})
		return
	}
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	personId, err := getPersonId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.ChatRoomService.DeleteChatRoom(personId, roomId)
	if errors.Is(err, model.NoRoleError) {
		c.JSON(http.StatusForbidden, gin.H{"error": "no role for this"})
		return
	}
	if err != nil {
		log.Printf("error deleting chat room: %s", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": "server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "success deleting chat room"})
}

func (h *Handler) exitFromChatRoom(c *gin.Context) {
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
	personId, err := getPersonId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.ChatRoomService.ExitFromChatRoom(personId, roomId)
	if err != nil {
		log.Printf("error exitng chat room: %s", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": "server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "success exiting chat room"})
}

func (h *Handler) deletePersonFromChatRoom(c *gin.Context) {
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
	personForDeletingStr, ok := c.GetQuery("person_id")
	if !ok {
		c.JSON(http.StatusBadRequest, "should be query room_id")
		return
	}
	personForDeleting, err := strconv.Atoi(personForDeletingStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	personId, err := getPersonId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.ChatRoomService.DeletePersonFromChatRoom(personId, personForDeleting, roomId)
	if errors.Is(err, model.NoRoleError) {
		c.JSON(http.StatusForbidden, gin.H{"error": "no role for this"})
		return
	}
	if err != nil {
		log.Printf("error deleting person from chat room: %s", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": "server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "success deleting person from chat room"})
}
