package entities

import (
	"time"
)

type Company struct {
	ID        string    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	CUIT      uint64    `gorm:"not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Branches  []Branch  `gorm:"foreignKey:CompanyID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type DeleteCompanyRequest struct {
	ID string `json:"id"`
}
