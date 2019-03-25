package models

import (
	"github.com/jinzhu/gorm"
	"../../lib/common"
)

type Post struct {
	gorm.Model
	Text   string `sql:"type:text;"`
	User   User   `gorm:"foreignkey:UserID"`
	UserID uint
}

func (p Post) Serialize() common.JSON {
	return common.JSON{
		"id":         p.ID,
		"text":       p.Text,
		"user":       p.User.Serialize(),
		"created_at": p.CreatedAt,
	}
}
