package postUser

import (
	"fmt"
	postUserDao "uttc-hackathon-backend/dao/postUser"
)

type UserUsecase struct {
	postUserDao *postUserDao.UserDAO
}

func NewUserUsecase(dao *postUserDao.UserDAO) *UserUsecase {
	return &UserUsecase{postUserDao: dao}
}

func (u *UserUsecase) RegisterUser(uid string, nickname string, sex string, birthyear int, birthdate int) (map[string]string, error) {
	if err := u.postUserDao.InsertUser(uid, nickname, sex, birthyear, birthdate); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	response := map[string]string{
		"message": "User registered successfully",
	}

	return response, nil
}
