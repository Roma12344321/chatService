package repository

import (
	"chatService/pkg/model"
	"github.com/jmoiron/sqlx"
)

type PersonRepository interface {
	CreatePerson(person *model.Person) (int, error)
	GetPerson(username, password string) (*model.Person, error)
	GetPersonById(id int) (model.Person, error)
	GetAllPersonByChatRoomId(roomId int) ([]model.PersonWithChatRoomRole, error)
	GetPersonByIdAndChatRoomId(personId, roomId int) (model.PersonWithChatRoomRole, error)
}

type ChatRoomRepository interface {
	CreateChatRoom(name string) (int, error)
	AddPersonToChatRoom(personId int, chatRoomId int, role string) error
	GetAllChatRoom(personId int) ([]model.ChatRoom, error)
	GetByPersonIdAndChatRoomId(personId, charRoomId int) (model.ChatRoom, error)
	DeletePersonFromChatRoom(personId, roomId int) error
	DeleteChatRoomById(roomId int) error
}

type MessageRepository interface {
	CreateMessage(text string, personId, chatRoomId int) (int, error)
	GetAllMessageForChatRoom(chatRoomId int) ([]model.Message, error)
	DeleteMessageById(id int) error
	GetMessageById(messageId int) (model.Message, error)
}

type Repository struct {
	PersonRepository
	ChatRoomRepository
	MessageRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		NewPersonRepositoryImpl(db),
		NewChatRoomRepositoryImpl(db),
		NewMessageRepositoryImpl(db),
	}
}
