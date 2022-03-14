package models

import (
	"future-fashion/dto"

	"gorm.io/gorm"

	"go.uber.org/zap"
)

type UserCRUDOperation interface {
	GetByID(uint) (*User, error)
	GetByUsername(string) (*User, error)
	GetAll() ([]*User, error)
	Insert(*User) (*User, error)
	Delete(uint) (*User, error)
	Update(uint, string, string, string, float32, float32, float32) (*User, error)
}

type UserCRUDOperationsImpl struct {
	DB     *gorm.DB
	Logger *zap.SugaredLogger
}

type User struct {
	gorm.Model
	Username string  `json:"username" gorm:"unique"`
	Password string  `json:"password"`
	DOB      string  `json:"dob"`
	Role     string  `json:"role"`
	Chest    float32 `json:"chest"`
	Waist    float32 `json:"waist"`
	Hip      float32 `json:"hip"`
	Orders   []Order
}

func (u *UserCRUDOperationsImpl) GetByID(id uint) (*User, error) {
	user := &User{}
	err := u.DB.First(user, id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserCRUDOperationsImpl) GetByUsername(username string) (*User, error) {
	user := &User{}
	err := u.DB.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserCRUDOperationsImpl) GetAll() ([]*User, error) {
	var users []*User
	err := u.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserCRUDOperationsImpl) Insert(userReq *dto.UserRequest) (*User, error) {
	user := &User{
		Username: userReq.Username,
		Password: userReq.Password,
		DOB:      userReq.DOB,
		Role:     userReq.Role,
	}

	err := u.DB.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserCRUDOperationsImpl) Delete(id uint) (*User, error) {
	foundUser, err := u.GetByID(id)
	if err != nil {
		return nil, err
	}
	//delete only if the user exists
	//permanently deleted with Unscoped().Delete()
	err = u.DB.Unscoped().Delete(foundUser, id).Error
	if err != nil {
		return nil, err
	}
	return foundUser, nil
}

func (u *UserCRUDOperationsImpl) Update(userReq *User) (*User, error) {
	foundUser, err := u.GetByID(userReq.ID)
	if err != nil {
		return nil, err
	}
	if userReq.Username != "" && foundUser.Username != userReq.Username {
		foundUser.Username = userReq.Username
	}
	if userReq.Password != "" && foundUser.Password != userReq.Password {
		foundUser.Password = userReq.Password
	}
	if userReq.DOB != "" && foundUser.DOB != userReq.DOB {
		foundUser.DOB = userReq.DOB
	}
	if userReq.Chest != 0 && foundUser.Chest != userReq.Chest {
		foundUser.Chest = userReq.Chest
	}
	if userReq.Waist != 0 && foundUser.Waist != userReq.Waist {
		foundUser.Waist = userReq.Waist
	}
	if userReq.Hip != 0 && foundUser.Hip != userReq.Hip {
		foundUser.Hip = userReq.Hip
	}

	//update user with all field
	err = u.DB.Save(foundUser).Error
	if err != nil {
		return nil, err
	}
	return foundUser, nil
}
