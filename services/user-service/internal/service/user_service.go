package service

import (
	"errors"
	"github.com/google/uuid"
	model "github.com/ppthinh/ChatApp/services/user-service/internal/models"
	pb "github.com/ppthinh/ChatApp/services/user-service/internal/proto/genproto"
	"github.com/ppthinh/ChatApp/services/user-service/internal/repository"
)

type userService struct {
	repo repository.UserRepository
}

type UserService interface {
	GetUserByEmail(email string) (*model.User, error)
	GetUserByPhoneNumber(phoneNumber string) (*model.User, error)
	CreateUser(req *pb.CreateUserRequest) (*model.User, error)
	UpdateUser(req *pb.UpdateUserRequest) (*model.User, error)
	DeleteUser(userID string) error

	GetUserForAuth(email string) (*model.User, error)
	// todo: implement sau
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

func (us *userService) CreateUser(req *pb.CreateUserRequest) (*model.User, error) {

	if user, _ := us.repo.GetUserByEmail(req.Email); user != nil {
		return nil, errors.New("email already exists")
	}
	if user, _ := us.repo.GetUserByPhoneNumber(req.PhoneNumber); user != nil {
		return nil, errors.New("phone number already exists")
	}

	user := &model.User{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	}
	if err := us.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) UpdateUser(req *pb.UpdateUserRequest) (*model.User, error) {
	user, err := us.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	user.Name = req.Name
	user.PhoneNumber = req.PhoneNumber
	user.Password = req.Password // Nên hash trong thực tế
	err = us.repo.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) DeleteUser(userID string) error {
	userUuid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("error converting string to UUID")
	}
	return us.repo.Delete(userUuid)
}

func (us *userService) GetUserForAuth(email string) (*model.User, error) {
	return us.repo.GetUserByEmail(email)
}
