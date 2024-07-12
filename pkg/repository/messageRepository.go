package repository

import (
	"chatService/pkg/model"
	"github.com/jmoiron/sqlx"
)

type MessageRepositoryImpl struct {
	db *sqlx.DB
}

func NewMessageRepositoryImpl(db *sqlx.DB) *MessageRepositoryImpl {
	return &MessageRepositoryImpl{db: db}
}

func (r *MessageRepositoryImpl) CreateMessage(text string, personId, chatRoomId int) (int, error) {
	query := `INSERT INTO message(text, person_id, chat_room_id) VALUES ($1,$2,$3)`
	row := r.db.QueryRow(query, text, personId, chatRoomId)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *MessageRepositoryImpl) GetAllMessageForChatRoom(chatRoomId int) ([]model.Message, error) {
	query := `SELECT * FROM message WHERE chat_room_id=$1`
	var res []model.Message
	err := r.db.Select(&res, query, chatRoomId)
	return res, err
}
