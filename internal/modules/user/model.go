package user

import "time"

// UserModel is DB model only (GORM).
type UserModel struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	Name      string `gorm:"size:120;not null"`
	Email     string `gorm:"size:160;uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserModel) TableName() string {
	return "users"
}
