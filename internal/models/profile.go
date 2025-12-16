package models

import (
	"database/sql"
	"time"
)

type Profile struct {
	ID        int64          `gorm:"column:id;primaryKey" json:"id"`
	UserID    int64          `gorm:"column:user_id;uniqueIndex;not null" json:"user_id"`
	Image     sql.NullString `gorm:"column:image;type:varchar(500)" json:"image"`
	Bio       sql.NullString `gorm:"column:bio;type:text" json:"bio"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;autoCreateTime;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;autoUpdateTime;not null" json:"updated_at"`
	User      *User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}
