package service

import (
	"github.com/google/uuid"
	model "github.com/ppthinh/ChatApp/services/user-service/internal/models"
	"github.com/ppthinh/ChatApp/services/user-service/internal/repository"
)

type userService struct {
	repo repository.UserRepository
}

type UserService interface {
	GetUserByEmail(email string) (*model.User, error)
	GetUserByPhoneNumber(phoneNumber string) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(userID uuid.UUID) error

	//GetFriendList (userID uuid.UUID) ([]model.User, error)
	//AddFriend (userId, friendId uuid.UUID) error
	//DeleteFriend (userId, friendId uuid.UUID) error
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (us *userService) GetUserByEmail(email string) (*model.User, error) {
	return us.repo.GetUserByEmail(email)
}

func (us *userService) GetUserByPhoneNumber(phoneNumber string) (*model.User, error) {
	return us.repo.GetUserByPhoneNumber(phoneNumber)
}

func (us *userService) CreateUser(user *model.User) error {
	return us.repo.Create(user)
}

func (us *userService) UpdateUser(user *model.User) error {
	return us.repo.Update(user)
}

func (us *userService) DeleteUser(userID uuid.UUID) error {
	return us.repo.Delete(userID)
}
