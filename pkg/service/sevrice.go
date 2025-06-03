package service

import (
	"github.com/Kiseshik/pet"
	"github.com/Kiseshik/pet/pkg/repository"
)

type Authorization interface {
	CreateUser(user pet.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list pet.TodoList) (int, error)
	GetAll(userId int) ([]pet.TodoList, error)
	GetById(userId, listId int) (pet.TodoList, error)
	Update(userId, listId int, input pet.UpdateListInput) error
	Delete(userId, listId int) error
}

type TodoItem interface {
	Create(userId, listId int, item pet.TodoItem) (int, error)
	GetAll(userId, listId int) ([]pet.TodoItem, error)
	GetById(userId, itemId int) (pet.TodoItem, error)
	Update(userId, itemId int, input pet.UpdateItemInput) error
	Delete(userId, itemId int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
