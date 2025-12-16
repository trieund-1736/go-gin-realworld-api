package models

import "time"

type Follow struct {
	ID         int64     `gorm:"column:id;primaryKey" json:"id"`
	FollowerID int64     `gorm:"column:follower_id;not null;uniqueIndex:idx_follow" json:"follower_id"`
	FolloweeID int64     `gorm:"column:followee_id;not null;uniqueIndex:idx_follow" json:"followee_id"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime;not null" json:"created_at"`
	Follower   *User     `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE" json:"-"`
	Followee   *User     `gorm:"foreignKey:FolloweeID;constraint:OnDelete:CASCADE" json:"-"`
}

func (Follow) TableName() string {
	return "follows"
}
