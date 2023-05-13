package model

import (
	"time"
)

type Comment struct {
	ID        uint       `json:"-" gorm:"primaryKey"`
	Uuid      string     `json:"uuid" gorm:"not null;uniqueIndex;"`
	ParentID  string     `json:"parentid"`
	Comment   string     `json:"comment" gorm:"not null"`
	Author    string     `json:"author" gorm:"not null"`
	UpdatedAt *time.Time `json:"update" gorm:"not null"`
	Favorite  bool       `json:"favorite" gorm:"not null;default:false"`
}
