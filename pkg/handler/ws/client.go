package ws

import (
	"bytes"
	"chatService/pkg/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
	postMethod     = "POST"
	deleteMethod   = "DELETE"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	roomId   int
	personId int
	service  *service.Service
}

type ReceivedMessage struct {
	Method  string `json:"method"`
	Content string `json:"content"`
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		rawMessage = bytes.TrimSpace(bytes.Replace(rawMessage, newline, space, -1))
		var resMessage ReceivedMessage
		err = json.Unmarshal(rawMessage, &resMessage)
		if err != nil {
			log.Printf("error unmarshaling rawMessage: %v", err)
			break
		}
		if resMessage.Method == postMethod {
			err = c.handlePostMethod(resMessage)
			if err != nil {
				break
			}
		} else if resMessage.Method == deleteMethod {
			err = c.handleDeleteMethod(resMessage)
			if err != nil {
				break
			}
		} else {
			break
		}
	}
}

func (c *Client) handlePostMethod(message ReceivedMessage) error {
	id, err := c.service.MessageService.CreateMessage(message.Content, c.personId, c.roomId)
	if err != nil {
		log.Printf("error saving message: %v", err)
		return err
	}
	c.hub.broadcast <- Message{Id: id, PersonId: c.personId, RoomId: c.roomId, Content: message.Content, Method: postMethod}
	return nil
}

func (c *Client) handleDeleteMethod(message ReceivedMessage) error {
	id, err := strconv.Atoi(message.Content)
	if err != nil {
		log.Printf("error read id: %v", err)
		return err
	}
	err = c.service.MessageService.DeleteMessageById(c.personId, c.roomId, id)
	if err != nil {
		log.Printf("error deleting message: %v", err)
		return err
	}
	c.hub.broadcast <- Message{Id: id, PersonId: c.personId, RoomId: c.roomId, Content: strconv.Itoa(id), Method: deleteMethod}
	return nil
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, personId, roomId int, service *service.Service) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), roomId: roomId, personId: personId, service: service}
	client.hub.register <- client
	go client.writePump()
	go client.readPump()
}
