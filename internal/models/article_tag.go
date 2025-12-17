package models

import "time"

type ArticleTag struct {
	ID        int64     `gorm:"column:id;primaryKey" json:"id"`
	ArticleID int64     `gorm:"column:article_id;not null;index;uniqueIndex:idx_article_tags" json:"article_id"`
	TagID     int64     `gorm:"column:tag_id;not null;index;uniqueIndex:idx_article_tags" json:"tag_id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;autoCreateTime;not null" json:"created_at"`
	Article   *Article  `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"-"`
	Tag       *Tag      `gorm:"foreignKey:TagID;constraint:OnDelete:CASCADE" json:"-"`
}
