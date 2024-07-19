package repository

import (
	"chatService/pkg/model"
	"github.com/jmoiron/sqlx"
	"time"
)

type MessageRepositoryImpl struct {
	db *sqlx.DB
}

func NewMessageRepositoryImpl(db *sqlx.DB) *MessageRepositoryImpl {
	return &MessageRepositoryImpl{db: db}
}

func (r *MessageRepositoryImpl) CreateMessage(message model.Message) (int, error) {
	query := `INSERT INTO message(text,date, person_id, chat_room_id) VALUES ($1,$2,$3,$4) RETURNING id`
	row := r.db.QueryRow(query, message.Text, message.Date, message.PersonId, message.ChatRoomId)
	var id int
	err := row.Scan(&id)
	return id, err
}

func (r *MessageRepositoryImpl) GetAllMessageForChatRoom(personId, chatRoomId int, date time.Time, limit int) ([]model.Message, error) {
	query := `SELECT * FROM message WHERE person_id=$1 AND chat_room_id=$2 AND date>$3 LIMIT $4`
	var res []model.Message
	err := r.db.Select(&res, query, personId, chatRoomId, date, limit)
	return res, err
}

func (r *MessageRepositoryImpl) DeleteMessageById(id int) error {
	query := `DELETE FROM message WHERE id=$1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *MessageRepositoryImpl) GetMessageById(messageId int) (model.Message, error) {
	query := `SELECT * FROM message WHERE id=$1`
	var res model.Message
	err := r.db.Get(&res, query, messageId)
	return res, err
}
