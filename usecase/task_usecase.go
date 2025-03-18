package usecase

import (
	"errors"
	"task-manager/domain"
)


type TaskUseCase struct {
	repo TaskRepository
}

func NewTaskUseCase(repo TaskRepository) *TaskUseCase {
	return &TaskUseCase{repo:repo}
}

func (uc *TaskUseCase) isValidStatus(status string) bool {
    return status == "pending" || status == "completed"
}

func (uc *TaskUseCase) Create(title, description, status *string) (domain.Task, error) {
	if !uc.isValidStatus(*status){
		return domain.Task{},errors.New("status must be 'pending' or 'completed'")
	}
	
	task:= domain.Task {
		Title: *title, 
		Description: *description,
		Status: *status,
	}
	return uc.repo.Create(task)
}

func (uc *TaskUseCase) GetAll() ([]domain.Task, error){
	return uc.repo.GetAll()
}

func (uc *TaskUseCase) GetByID(id int)(domain.Task, error) {
	return uc.repo.GetByID(id)
}

func (uc *TaskUseCase) Update(id int, title, description, status *string) (domain.Task, error) {
	
   task, err:= uc.repo.GetByID(id)
   if err != nil {
	return domain.Task{}, err
   }
if title != nil {
	task.Title = *title
}
if description != nil {
	task.Description = *description
}
if status != nil {
	if !uc.isValidStatus(*status){
		return domain.Task{}, errors.New("status must be 'pending' or 'completed'")
	}
	task.Status = *status
}
	
	return uc.repo.Update(task)
}

func (uc *TaskUseCase) Delete(id int) error {
	return uc.repo.Delete(id)
}
