package entities

import (
	"time"
)

type Branch struct {
	ID             string    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name           string    `gorm:"type:varchar(255);not null"`
	Address        string    `gorm:"type:varchar(255);not null"`
	OpenTime       time.Time `gorm:"type:time;not null"`
	CloseTime      time.Time `gorm:"type:time;not null"`
	DayOfWeekOpen  string    `gorm:"type:varchar(50);not null"`
	CompanyID      string    `gorm:"type:uuid;not null"`
	Company        Company   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	HeaderImgPath  string    `gorm:"type:varchar(255);not null"`
	ProfileImgPath string    `gorm:"type:varchar(255);not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

type DeleteBranchRequest struct {
	ID string `json:"id"`
}
