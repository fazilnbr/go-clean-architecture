package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(users User) error
	Login(userName, Password string) (User, error)
	UpdateUser(user User) error
	DeleteUser(id uint) error
}

type service struct {
	dao DAO
}

func NewService(dao DAO) Service {
	return service{
		dao: dao,
	}
}

func (s service) CreateUser(users User) error {
	hashedPassword, err := s.hashPassword(users.Password)
	if err != nil {
		return err
	}
	users.Password = hashedPassword
	err = s.dao.CreateUser(users)
	return err
}

func (s service) Login(userName, Password string) (User, error) {
	user, err := s.dao.GetUser(userName)
	if err != nil {
		return user, err
	}
	valid := s.compareHashAndPassword(user.Password, Password)
	if !valid {
		return user, errors.New("user name or password is invalid")
	}
	return user, nil
}

func (s service) UpdateUser(user User) error {
	if user.Password != "" {
		hashedPassword, err := s.hashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}
	err := s.dao.UpdateUser(user)
	return err
}

func (s service) DeleteUser(id uint) error {
	err := s.dao.DeleteUser(id)
	return err
}

func (s service) hashPassword(password string) (string, error) {
	// Hash the password using the default cost
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return string(hashedPassword), err
	}
	return string(hashedPassword), err

}

func (s service) compareHashAndPassword(hashedPassword, password string) bool {
	// Compare the hashed password with the plaintext password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
