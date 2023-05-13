package model

import (
	"auth-service/auth/pkg"

	"gorm.io/gorm"
)

type User struct {
	Base
	Password   string `json:"password" gorm:"varchar(100);not null"`
	IsVerified bool   `json:"is_verified"`
	Email      string `json:"email" gorm:"varchar(100);not null;uniqueIndex"`
	Username   string `json:"username" gorm:"varchar(100);"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	id, err := pkg.GenerateSnowflake()
	if err != nil {
		return err
	}
	u.ID = id.Int64()
	return nil
}
