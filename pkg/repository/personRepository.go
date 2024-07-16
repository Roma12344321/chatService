package repository

import (
	"chatService/pkg/model"
	"github.com/jmoiron/sqlx"
)

type PersonRepositoryImpl struct {
	db *sqlx.DB
}

func NewPersonRepositoryImpl(db *sqlx.DB) *PersonRepositoryImpl {
	return &PersonRepositoryImpl{db: db}
}

func (r *PersonRepositoryImpl) CreatePerson(person *model.Person) (int, error) {
	query := "INSERT INTO person(username, password,role) VALUES ($1,$2,$3) RETURNING id"
	row := r.db.QueryRow(query, person.Username, person.Password, person.Role)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PersonRepositoryImpl) GetPerson(username, password string) (*model.Person, error) {
	var person model.Person
	query := "SELECT * FROM person WHERE username=$1 AND password=$2"
	if err := r.db.Get(&person, query, username, password); err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *PersonRepositoryImpl) GetPersonById(id int) (model.Person, error) {
	var person model.Person
	query := "SELECT * FROM person WHERE id=$1"
	if err := r.db.Get(&person, query, id); err != nil {
		return person, err
	}
	return person, nil
}

func (r *PersonRepositoryImpl) GetAllPersonByChatRoomId(roomId int) ([]model.PersonWithChatRoomRole, error) {
	query := `SELECT person.id as "person.id",username as "person.username",password as "person.password",person.role 
    as "person.role",person_chat_room.role FROM person JOIN person_chat_room ON person.id = person_chat_room.person_id WHERE chat_room_id=$1`
	var res []model.PersonWithChatRoomRole
	err := r.db.Select(&res, query, roomId)
	return res, err
}

func (r *PersonRepositoryImpl) GetPersonByIdAndChatRoomId(personId, roomId int) (model.PersonWithChatRoomRole, error) {
	query := `SELECT person.id as "person.id",username as "person.username",password as "person.password",person.role 
    as "person.role",person_chat_room.role FROM person JOIN person_chat_room ON person.id = person_chat_room.person_id WHERE person_id=$1 AND chat_room_id=$2`
	var res model.PersonWithChatRoomRole
	err := r.db.Get(&res, query, personId, roomId)
	return res, err
}
