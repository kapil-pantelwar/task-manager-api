package usecase

import "task-manager/domain"


type TaskRepository interface {
	Create(task domain.Task) (domain.Task, error)
	GetAll() ([]domain.Task, error)
	GetByID(id int) (domain.Task, error)
	Update(task domain.Task) (domain.Task, error)
	Delete(id int) error
}


