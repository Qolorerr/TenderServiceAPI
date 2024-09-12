package models

import (
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	VersionID   uuid.UUID `json:"versionId" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Description string    `json:"description" gorm:"type:text;not null"`
	Status      string    `json:"status" gorm:"type:varchar(20);default:'Created';not null"`
	TenderId    uuid.UUID `json:"tenderId" gorm:"not null"`
	AuthorType  string    `json:"authorType" gorm:"type:varchar(20);not null"`
	AuthorId    uuid.UUID `json:"authorId" gorm:"not null"`
	Author      Employee  `json:"author" gorm:"foreignkey:AuthorId;references:id"`
	Version     int32     `json:"version" gorm:"default:1;not null"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

type BidEdit struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
