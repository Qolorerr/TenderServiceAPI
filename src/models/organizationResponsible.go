package models

import (
	"github.com/google/uuid"
)

type OrganizationResponsible struct {
	ID             uuid.UUID    `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OrganizationID uuid.UUID    `json:"organizationId" gorm:"not null"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	UserID         uuid.UUID    `json:"userId" gorm:"not null"`
	Employee       Employee     `gorm:"foreignKey:UserID"`
}

func (OrganizationResponsible) TableName() string {
	return "organization_responsible"
}
