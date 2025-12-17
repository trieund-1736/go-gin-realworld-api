package models

import "time"

type Favorite struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"id"`
	UserID    int64     `gorm:"column:user_id;not null;index;uniqueIndex:idx_favorites" json:"user_id"`
	ArticleID int64     `gorm:"column:article_id;not null;index;uniqueIndex:idx_favorites" json:"article_id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime;not null" json:"created_at"`
	User      *User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Article   *Article  `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"-"`
}
