package service

import (
	"chatService/pkg/model"
	"chatService/pkg/repository"
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
