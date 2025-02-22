package repository

import (
	"errors"
	"github.com/google/uuid"
	model "github.com/ppthinh/ChatApp/services/user-service/internal/models"
	"gorm.io/gorm"
)

var (
	ErrorUserNotFound   = errors.New("User not found")
	ErrEmailExist       = errors.New("Email already exist")
	ErrPhoneNumberExist = errors.New("Phone number already exist")
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	GetUserByEmail(email string) (*model.User, error)
	GetUserByPhoneNumber(phoneNumber string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(userID uuid.UUID) error

	GetFriendList(userID uuid.UUID) ([]model.User, error)
	AddFriend(userId, friendId uuid.UUID) error
	DeleteFriend(userId uuid.UUID, friendId uuid.UUID) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// -----user
func (ur *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := ur.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) GetUserByPhoneNumber(phoneNumber string) (*model.User, error) {
	var user model.User
	err := ur.db.Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) Create(user *model.User) error {
	var existingUser model.User
	err := ur.db.Create(user).Error
	if err != nil {
		if err := ur.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			return ErrEmailExist
		}

		if err := ur.db.Where("phone_number = ?", user.PhoneNumber).First(&existingUser).Error; err == nil {
			return ErrPhoneNumberExist
		}

		return err

	}
	return nil
}
func (ur *userRepository) Update(user *model.User) error {
	var existingUser model.User
	err := ur.db.Where("email = ? AND id <> ?", user.Email, user.ID).First(&existingUser).Error
	if err == nil {
		return ErrEmailExist
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	err = ur.db.Where("phone_number = ? AND id <> ?", user.PhoneNumber, user.ID).First(&existingUser).Error
	if err == nil {
		return ErrPhoneNumberExist
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return ur.db.Model(user).Updates(map[string]interface{}{
		"name":         user.Name,
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
	}).Error
}
func (ur *userRepository) Delete(userID uuid.UUID) error {
	return ur.db.Transaction(func(tx *gorm.DB) error {
		// xoa user trong danh sach ban be cua user & user khac
		err := tx.Where("user_id = ? or friend_id = ?", userID, userID).Delete(&model.FriendShip{}).Error
		if err != nil {
			return err
		}

		// xoa user trong danh sach yeu cau ket ban cua user & user khac
		err = tx.Where("from_user_id = ? or to_user_id = ?", userID, userID).Delete(&model.FriendRequest{}).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&model.User{}, userID).Error
		if err != nil {
			return err
		}

		return nil
	})
}

// -----friendship
func (ur *userRepository) GetFriendList(userID uuid.UUID) ([]model.User, error) {
	var friends []model.User
	// SELECT users.id, users.name
	// FROM friendships
	// JOIN users ON friendships.friend_id = users.id
	// WHERE friendships.user_id = 'uuid';
	err := ur.db.Table("friendship").
		Select("users.id, users.name").
		Joins("JOIN users ON friendship.friend_id = users.id").
		Where("friendship.user_id = ?", userID).
		Find(&friends).Error
	return friends, err
}

func (ur *userRepository) AddFriend(UserID, FriendID uuid.UUID) error {
	friendShip := model.FriendShip{
		UserID:   UserID,
		FriendID: FriendID,
	}

	return ur.db.Create(&friendShip).Error
}

func (ur *userRepository) DeleteFriend(userID, friendID uuid.UUID) error {
	return ur.db.Where("user_id = ? and friend_id = ?", userID, friendID).Delete(&model.FriendShip{}).Error
}
