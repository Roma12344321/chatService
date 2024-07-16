package model

type Message struct {
	Id         int    `json:"id"`
	Text       string `json:"text"`
	PersonId   int    `json:"person_id" db:"person_id"`
	ChatRoomId int    `json:"chat_room_id" db:"chat_room_id"`
}
