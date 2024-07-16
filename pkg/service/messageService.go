package service

import (
	"chatService/pkg/model"
	"chatService/pkg/repository"
	"errors"
)

type MessageServiceImpl struct {
	repo *repository.Repository
}

func NewMessageServiceImpl(repo *repository.Repository) *MessageServiceImpl {
	return &MessageServiceImpl{repo: repo}
}

func (s *MessageServiceImpl) CreateMessage(text string, personId, chatRoomId int) (int, error) {
	return s.repo.MessageRepository.CreateMessage(text, personId, chatRoomId)
}

func (s *MessageServiceImpl) GetAllMessageForChatRoom(chatRoomId int) ([]model.Message, error) {
	return s.repo.MessageRepository.GetAllMessageForChatRoom(chatRoomId)
}

func (s *MessageServiceImpl) DeleteMessageById(personId, roomId, messageId int) error {
	person, err := s.repo.PersonRepository.GetPersonByIdAndChatRoomId(personId, roomId)
	if err != nil {
		return err
	}
	if person.Person == nil {
		return errors.New("person with id not found")
	}
	message, err := s.repo.MessageRepository.GetMessageById(messageId)
	if err != nil {
		return err
	}
	if person.Person.Role == model.RoleAdmin || person.RoomRole == model.RoleAdmin || person.Person.Id == message.PersonId {
		return s.repo.MessageRepository.DeleteMessageById(messageId)
	}
	return errors.New("no role for this")
}
