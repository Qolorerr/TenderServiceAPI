package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Description string    `json:"description"`
	Type        string    `json:"type" gorm:"type:organization_type"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	// Связь с ответственными сотрудниками
	ResponsibleEmployees []OrganizationResponsible `gorm:"foreignKey:OrganizationID"`
}

func (Organization) TableName() string {
	return "organization"
}
