package userService

import "gorm.io/gorm"

type UserRepository interface {
	CreateUser(user User) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUserById(id uint, user User) (User, error)
	DeleteUserById(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user User) (User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *userRepository) GetAllUsers() ([]User, error) {
	var users []User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) UpdateUserById(id uint, user User) (User, error) {
	var existing User
	if err := r.db.First(&existing, id).Error; err != nil {
		return User{}, err
	}

	existing.Email = user.Email
	existing.Password = user.Password

	if err := r.db.Save(&existing).Error; err != nil {
		return User{}, err
	}
	return existing, nil
}

func (r *userRepository) DeleteUserById(id uint) error {
	if err := r.db.Delete(&User{}, id).Error; err != nil {
		return err
	}
	return nil
}
