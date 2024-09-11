package models

import (
	"time"

	"github.com/google/uuid"
)

type Tender struct {
	ID              uuid.UUID    `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	VersionID       uuid.UUID    `json:"versionId" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name            string       `json:"name" gorm:"type:varchar(100);not null"`
	Description     string       `json:"description" gorm:"type:text;not null"`
	ServiceType     string       `json:"serviceType" gorm:"type:varchar(20);not null"`
	Status          string       `json:"status" gorm:"type:varchar(20);default:'Created';not null"`
	OrganizationId  string       `json:"organizationId" gorm:"type:varchar(20);not null"`
	Organization    Organization `json:"organization" gorm:"foreignkey:OrganizationId;references:id"`
	CreatorUsername string       `json:"creatorUsername" gorm:"-"`
	Version         int32        `json:"version" gorm:"default:1;not null"`
	CreatedAt       time.Time    `json:"createdAt" gorm:"autoCreateTime"`
}

type TenderEdit struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"serviceType"`
}
