package service

import (
	"chatService/pkg/model"
	"chatService/pkg/repository"
	"errors"
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
	err = s.repo.AddPersonToChatRoom(fromPersonId, roomId, model.RoleAdmin)
	if err != nil {
		return err
	}
	err = s.repo.ChatRoomRepository.AddPersonToChatRoom(toPersonId, roomId, model.RoleUser)
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

func (s *ChatRoomServiceImpl) GetAllPersonByChatRoomId(roomId int) ([]model.PersonWithChatRoomRole, error) {
	return s.repo.PersonRepository.GetAllPersonByChatRoomId(roomId)
}

func (s *ChatRoomServiceImpl) DeleteChatRoom(personId, roomId int) (bool, error) {
	ok, err := s.isPersonHasRole(personId, roomId)
	if err != nil {
		return false, err
	}
	if ok {
		err := s.repo.ChatRoomRepository.DeleteChatRoomById(roomId)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (s *ChatRoomServiceImpl) ExitFromChatRoom(personId, roomId int) error {
	if err := s.repo.ChatRoomRepository.DeletePersonFromChatRoom(personId, roomId); err != nil {
		return err
	}
	people, err := s.repo.PersonRepository.GetAllPersonByChatRoomId(roomId)
	if err != nil {
		return err
	}
	if len(people) == 0 {
		return s.repo.ChatRoomRepository.DeleteChatRoomById(roomId)
	}
	return nil
}

func (s *ChatRoomServiceImpl) DeletePersonFromChatRoom(personId, personForDeletingId, roomId int) error {
	ok, err := s.isPersonHasRole(personId, roomId)
	if err != nil {
		return err
	}
	if ok {
		return s.repo.DeletePersonFromChatRoom(personForDeletingId, roomId)
	}
	return errors.New("no role")
}

func (s *ChatRoomServiceImpl) isPersonHasRole(personId, roomId int) (bool, error) {
	people, err := s.repo.PersonRepository.GetAllPersonByChatRoomId(roomId)
	if err != nil {
		return false, err
	}
	var personThatDeleting model.PersonWithChatRoomRole
	for _, person := range people {
		if person.Person.Id == personId {
			personThatDeleting = person
		}
	}
	if personThatDeleting.Person != nil {
		if personThatDeleting.Person.Role == model.RoleAdmin || personThatDeleting.RoomRole == model.RoleAdmin || len(people) == 1 {
			return true, nil
		}
	}
	return false, nil
}
