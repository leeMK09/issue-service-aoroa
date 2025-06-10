package infrastructure

import userModel "issue-service-aoroa/user/model"

type UserRepository interface {
	GetByID(id uint) (*userModel.User, error)
	GetAll() []userModel.User
}

type userRepository struct {
	users []userModel.User
}

func NewUserRepository() UserRepository {
	return &userRepository{
		users: []userModel.User{
			{ID: 1, Name: "김개발"},
			{ID: 2, Name: "이디자인"},
			{ID: 3, Name: "박기획"},
		},
	}
}

func (r *userRepository) GetByID(id uint) (*userModel.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, nil
}

func (r *userRepository) GetAll() []userModel.User {
	return r.users
}