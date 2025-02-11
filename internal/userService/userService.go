package userService

import "errors"

type UserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserService struct {
	repo UserRepository
}

var ErrUserNotFound = errors.New("user not found")

func NewService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user User) (User, error) {
	return s.repo.CreateUser(user)
}

func (s *UserService) GetAllUsers() ([]User, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) UpdateUserById(id uint, user User) (User, error) {
	return s.repo.UpdateUserById(id, user)
}

func (s *UserService) DeleteUserById(id uint) error {
	users, err := s.GetAllUsers()
	if err != nil {
		return err
	}

	var userToDelete *User
	for _, user := range users {
		if user.ID == id {
			userToDelete = &user
			break
		}
	}

	if userToDelete == nil {
		return ErrUserNotFound
	}

	if err := s.repo.DeleteUserById(id); err != nil {
		return err
	}
	return nil
}
