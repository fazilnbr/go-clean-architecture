package users

import (
	"gorm.io/gorm"
)

type DAO interface {
	CreateUser(users User) error
	GetUser(userName string) (User, error)
	UpdateUser(user User) error
	DeleteUser(id uint) error
}

type sqlDao struct {
	db *gorm.DB
}

func NewDAO(db *gorm.DB) DAO {
	return sqlDao{db: db}
}

func (d sqlDao) CreateUser(user User) error {
	err := d.db.Create(&user).Error
	return err
}

func (d sqlDao) GetUser(userName string) (User, error) {
	var user User
	err := d.db.Where(User{UserName: userName}).First(&user).Error
	return user, err
}

func (d sqlDao) UpdateUser(user User) error {
	err := d.db.Model(&User{}).Where(User{ID: user.ID}).Updates(user).Error
	return err
}

func (d sqlDao) DeleteUser(id uint) error {
	// delete permanently
	err := d.db.Unscoped().Where(User{ID: id}).Delete(&User{}).Error
	return err
}
