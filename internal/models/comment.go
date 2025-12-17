package models

import "time"

type Comment struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"id"`
	Body      string    `gorm:"column:body;type:text;not null" json:"body"`
	ArticleID int64     `gorm:"column:article_id;not null;index" json:"article_id"`
	AuthorID  int64     `gorm:"column:author_id;not null;index" json:"author_id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;autoUpdateTime;not null" json:"updated_at"`
	Article   *Article  `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"-"`
	Author    *User     `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE" json:"-"`
}
