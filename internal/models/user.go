package models

import "time"

type User struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"id"`
	Username  string    `gorm:"column:username;type:varchar(255);uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"column:email;type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"column:password;type:varchar(255);not null" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime;not null" json:"updated_at"`
	Profile   *Profile  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}
