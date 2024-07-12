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
