package user

import (
	"quiz-builder/entity"
	"quiz-builder/logger"

	"quiz-builder/repository/quiz_builder"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type IUserService interface {
	WithTrx(trxHandle *gorm.DB) IUserService
	GetOneUser(uuid uuid.UUID) (entity.User, error)
	GetAllUser() ([]entity.User, error)
	CreateUser(entity.User) error
	UpdateUser(entity.User) error
	DeleteUser(uuid uuid.UUID) error
}

// UserService service layer
type UserService struct {
	logger     logger.Logger
	repository quiz_builder.QuizBuilderRepository
}

// NewUserService creates a new userservice
func NewUserService(
	logger logger.Logger,
	repository quiz_builder.QuizBuilderRepository,
) UserService {
	return UserService{
		logger:     logger,
		repository: repository,
	}
}

// WithTrx delegates transaction to repository database
func (s UserService) WithTrx(trxHandle *gorm.DB) UserService {
	s.repository = s.repository.WithTrx(trxHandle)
	return s
}

// GetOneUser gets one user
func (s UserService) GetOneUser(uuid uuid.UUID) (user entity.User, err error) {
	return user, s.repository.Find(&user, "UUID = ?", uuid.String()).Error
}

// GetOneUserByID gets one user by id
func (s UserService) GetOneUserByID(userID uint) (user entity.User, err error) {
	return user, s.repository.Find(&user, userID).Error
}

// GetOneUserByEmail gets one user by email
func (s UserService) GetOneUserByEmail(email string) (user entity.User, err error) {
	return user, s.repository.Find(&user, "Email = ?", email).Error
}

// GetAllUser get all the user
func (s UserService) GetAllUser() (users []entity.User, err error) {
	return users, s.repository.Find(&users).Error
}

// CreateUser call to create the user
func (s UserService) CreateUser(user entity.User) error {
	return s.repository.Create(&user).Error
}

// UpdateUser updates the user
func (s UserService) UpdateUser(user entity.User) error {
	return s.repository.Save(&user).Error
}

// DeleteUser deletes the user
func (s UserService) DeleteUser(uuid uuid.UUID) error {
	return s.repository.Delete(&entity.User{}, "UUID = ?", uuid.String()).Error
}
