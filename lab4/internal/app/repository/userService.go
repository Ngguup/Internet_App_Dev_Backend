package repository

import (
	"errors"
	"lab1/internal/app/ds"
	"lab1/internal/app/role"
)

func (r *Repository) CreateUser(user *ds.Users) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetUserByID(id uint) (*ds.Users, error) {
	var user ds.Users
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(id uint, login, password string, role role.Role) error {
	var user ds.Users
	if err := r.db.First(&user, id).Error; err != nil {
		return err
	}

	if login != "" {
		user.Login = login
	}

	if password != "" {
		user.Password = password
	}

	user.Role = role //!!!

	if err := r.db.Save(&user).Error; err != nil {
		return errors.New("failed to update user")
	}
	return nil
}









func (r *Repository) Register(user *ds.Users) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetUserByLogin(login string) (*ds.Users, error) {
	user := &ds.Users{}

	err := r.db.Where("login = ?", login).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

