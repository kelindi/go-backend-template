package main

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Username     string    `gorm:"not null;unique"`
	Password     string    `gorm:"not null"`
	Email string    `gorm:"not null;unique"`
}