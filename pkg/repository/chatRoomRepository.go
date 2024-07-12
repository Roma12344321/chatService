package repository

import (
	"chatService/pkg/model"
	"github.com/jmoiron/sqlx"
	"time"
)

type ChatRoomRepositoryImpl struct {
	db *sqlx.DB
}

func NewChatRoomRepositoryImpl(db *sqlx.DB) *ChatRoomRepositoryImpl {
	return &ChatRoomRepositoryImpl{db: db}
}

func (r *ChatRoomRepositoryImpl) CreateChatRoom(name string) (int, error) {
	query := `INSERT INTO chat_room(name, date) VALUES ($1,$2) RETURNING id`
	row := r.db.QueryRow(query, name, time.Now())
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ChatRoomRepositoryImpl) AddPersonToChatRoom(personId int, chatRoomId int) error {
	query := `INSERT INTO person_chat_room(person_id, chat_room_id) VALUES ($1,$2)`
	if _, err := r.db.Exec(query, personId, chatRoomId); err != nil {
		return err
	}
	return nil
}

func (r *ChatRoomRepositoryImpl) GetAllChatRoom(personId int) ([]model.ChatRoom, error) {
	query := `SELECT chat_room.id,chat_room.name,chat_room.date FROM chat_room LEFT JOIN 
    person_chat_room ON chat_room.id=person_chat_room.chat_room_id WHERE person_id=$1`
	var res []model.ChatRoom
	err := r.db.Select(&res, query, personId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *ChatRoomRepositoryImpl) GetByPersonIdAndChatRoomId(personId, charRoomId int) (model.ChatRoom, error) {
	query := `SELECT id,name,date FROM person_chat_room JOIN chat_room ON person_chat_room.chat_room_id = chat_room.id WHERE person_id=$1 AND chat_room_id=$2`
	var res model.ChatRoom
	err := r.db.Get(&res, query, personId, charRoomId)
	return res, err
}
