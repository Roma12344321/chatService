package service

import (
	"chatService/pkg/model"
	"chatService/pkg/repository"
)

type ChatRoomServiceImpl struct {
	repo *repository.Repository
}

func NewChatRoomServiceImpl(repo *repository.Repository) *ChatRoomServiceImpl {
	return &ChatRoomServiceImpl{repo: repo}
}

func (s *ChatRoomServiceImpl) CreateChatRoom(fromPersonId, toPersonId int) error {
	toPerson, err := s.repo.PersonRepository.GetPersonById(toPersonId)
	if err != nil {
		return err
	}
	roomId, err := s.repo.ChatRoomRepository.CreateChatRoom(toPerson.Username)
	if err != nil {
		return err
	}
	err = s.repo.AddPersonToChatRoom(fromPersonId, roomId)
	if err != nil {
		return err
	}
	err = s.repo.ChatRoomRepository.AddPersonToChatRoom(toPersonId, roomId)
	if err != nil {
		return err
	}
	return nil
}

func (s *ChatRoomServiceImpl) GetAllChatRoom(personId int) ([]model.ChatRoom, error) {
	return s.repo.ChatRoomRepository.GetAllChatRoom(personId)
}

func (s *ChatRoomServiceImpl) GetByPersonIdAndChatRoomId(personId, charRoomId int) (model.ChatRoom, error) {
	return s.repo.ChatRoomRepository.GetByPersonIdAndChatRoomId(personId, charRoomId)
}
