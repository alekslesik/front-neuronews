package service

import "github.com/alekslesik/neuro-news/internal/app/model"

type UserService interface {
	GetUserByID(id int) (*model.User, error)
}

type userService struct {
	ur model.UserModel
}

func NewUserService(userRepository model.UserModel) UserService {
	return &userService{
		ur: userRepository,
	}
}

func (us *userService) GetUserByID(id int) (*model.User, error) {
	return us.ur.GetUserByID(id)
}
