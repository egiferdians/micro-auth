package models

import (
	"github.com/egiferdians/micro/util/validator"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	IDUser        	uuid.UUID `gorm:"type:uuid;primary_key;" json:"id_user" db:"id_user"`
	Fullname     	string    `gorm:"size:255;not null;" json:"fullname" validate:"required"`
	Email     		string    `gorm:"size:255;not null;" json:"email" validate:"required,email"`
	Password  		string    `gorm:"size:255;not null;" json:"password" validate:"required"`
}

func (m *User) Validate() error {
	err := validator.V([]interface{}{
		m.Email,
		m.Password,
	}, m)
	if err != nil {
		return err
	}
	return nil
}
