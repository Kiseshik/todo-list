package repository

import (
	"github.com/Kiseshik/pet"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user pet.User) (int, error)
	GetUser(username, password string) (pet.User, error)
}

type TodoList interface {
	Create(userId int, list pet.TodoList) (int, error)
	GetAll(userId int) ([]pet.TodoList, error)
	GetById(userId, listId int) (pet.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input pet.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, item pet.TodoItem) (int, error)
	GetAll(userId, listId int) ([]pet.TodoItem, error)
	GetById(userId, itemId int) (pet.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input pet.UpdateItemInput) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
