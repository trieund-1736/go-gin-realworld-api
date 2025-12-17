package models

import "time"

type Article struct {
	ID             int64         `gorm:"column:id;primaryKey" json:"id"`
	Slug           string        `gorm:"column:slug;type:varchar(500);uniqueIndex;not null" json:"slug"`
	Title          string        `gorm:"column:title;type:varchar(500);not null" json:"title"`
	Description    string        `gorm:"column:description;type:text;not null" json:"description"`
	Body           string        `gorm:"column:body;type:text;not null" json:"body"`
	AuthorID       int64         `gorm:"column:author_id;not null;index" json:"author_id"`
	FavoritesCount int           `gorm:"column:favorites_count;default:0" json:"favorites_count"`
	CreatedAt      time.Time     `gorm:"column:created_at;type:timestamp;autoCreateTime;not null" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"column:updated_at;type:timestamp;autoUpdateTime;not null" json:"updated_at"`
	Author         *User         `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE" json:"-"`
	Comments       []*Comment    `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"-"`
	ArticleTags    []*ArticleTag `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"-"`
	Favorites      []*Favorite   `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"-"`
}
