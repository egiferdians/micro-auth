package mysql

import (
	"context"

	"github.com/egiferdians/micro-auth/auth"
	_interface "github.com/egiferdians/micro-auth/entity/interface"
	"github.com/egiferdians/micro-auth/models"
	"github.com/jinzhu/gorm"
)

type repository struct {
	DB *gorm.DB
}

func NewDBReadWriter(db *gorm.DB) _interface.Repository {
	return &repository{
		DB: db,
	}
}

func (r *repository) Login(_ context.Context, usr models.User) (*models.User, *auth.Authenticated, error) {
	var err error
	err = r.DB.Table("man_access.user").Where(
		"email = ?",
		usr.Email,
	).First(
		&usr,
	).Error
	if err != nil {
		return nil, nil, err
	}
	return &usr, &auth.Authenticated{}, nil
}