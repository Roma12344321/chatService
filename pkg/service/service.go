package service

import (
	"chatService/pkg/model"
	"chatService/pkg/repository"
)

type AuthService interface {
	Registration(person *model.Person) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	GetPersonById(id int) (model.Person, error)
}

type ChatRoomService interface {
	CreateChatRoom(fromPersonId, toPersonId int) error
	GetAllChatRoom(personId int) ([]model.ChatRoom, error)
	GetByPersonIdAndChatRoomId(personId, charRoomId int) (model.ChatRoom, error)
	GetAllPersonByChatRoomId(roomId int) ([]model.PersonWithChatRoomRole, error)
	DeleteChatRoom(personId, roomId int) (bool, error)
	ExitFromChatRoom(personId, roomId int) error
	DeletePersonFromChatRoom(personId, personForDeletingId, roomId int) error
}

type MessageService interface {
	CreateMessage(text string, personId, chatRoomId int) (int, error)
	GetAllMessageForChatRoom(chatRoomId int) ([]model.Message, error)
	DeleteMessageById(personId, roomId, messageId int) error
}

type Service struct {
	AuthService
	ChatRoomService
	MessageService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		NewAuthServiceImpl(repository),
		NewChatRoomServiceImpl(repository),
		NewMessageServiceImpl(repository),
	}
}
