package models

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username  string    `json:"username" gorm:"type:varchar(50);unique;not null"`
	FirstName string    `json:"firstName" gorm:"type:varchar(50)"`
	LastName  string    `json:"lastName" gorm:"type:varchar(50)"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	// Связь с ответственностью в организации
	OrganizationsResponsible []OrganizationResponsible `gorm:"foreignKey:UserID"`
}

func (Employee) TableName() string {
	return "employee"
}
