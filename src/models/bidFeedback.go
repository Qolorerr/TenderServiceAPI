package models

import (
	"github.com/google/uuid"
	"time"
)

type BidFeedback struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Description string    `json:"description" gorm:"type:text;not null"`
	BidId       uuid.UUID `json:"authorId" gorm:"not null"`
	Bid         Bid       `json:"author" gorm:"foreignkey:AuthorId;references:id"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
}
