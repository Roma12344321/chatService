package service

import (
	"chatService/pkg/model"
	"chatService/pkg/repository"
	"errors"
	"time"
)

type MessageServiceImpl struct {
	repo *repository.Repository
}

func NewMessageServiceImpl(repo *repository.Repository) *MessageServiceImpl {
	return &MessageServiceImpl{repo: repo}
}

func (s *MessageServiceImpl) CreateMessage(message model.Message) (int, error) {
	return s.repo.MessageRepository.CreateMessage(message)
}

func (s *MessageServiceImpl) GetAllMessageForChatRoom(personId, chatRoomId int, date time.Time, limit int) ([]model.Message, error) {
	return s.repo.MessageRepository.GetAllMessageForChatRoom(personId, chatRoomId, date, limit)
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
