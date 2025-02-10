package taskService

import (
	"errors"
)

type TaskResponse struct {
	ID     uint   `json:"id"`
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
}

type TaskService struct {
	repo TaskRepository
}

var ErrTaskNotFound = errors.New("task not found")

func NewService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(task Task) (Task, error) {
	return s.repo.CreateTask(task)
}

func (s *TaskService) GetAllTasks() ([]Task, error) {
	return s.repo.GetAllTasks()
}

func (s *TaskService) UpdateTaskById(id uint, task Task) (Task, error) {
	return s.repo.UpdateTaskByID(id, task)
}

func (s *TaskService) DeleteTaskById(id uint) error {
	tasks, err := s.GetAllTasks()
	if err != nil {
		return err
	}

	var taskToDelete *Task
	for _, task := range tasks {
		if task.ID == id {
			taskToDelete = &task
			break
		}
	}

	if taskToDelete == nil {
		return ErrTaskNotFound
	}

	if err := s.repo.DeleteTaskByID(id); err != nil {
		return err
	}
	return nil
}
