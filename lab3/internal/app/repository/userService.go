package repository

import (
	"errors"
	"lab1/internal/app/ds"
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

func (r *Repository) UpdateUser(id uint, login, password string, isModerator *bool) error {
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

	if isModerator != nil {
		user.IsModerator = *isModerator
	}

	if err := r.db.Save(&user).Error; err != nil {
		return errors.New("failed to update user")
	}
	return nil
}
