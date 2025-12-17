package models

import "time"

type Tag struct {
	ID        int64      `gorm:"column:id;primaryKey" json:"id"`
	Name      string     `gorm:"column:name;type:varchar(255);uniqueIndex;not null" json:"name"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp;autoCreateTime;not null" json:"created_at"`
	Articles  []*Article `gorm:"many2many:article_tags;foreignKey:ID;joinForeignKey:TagID;references:ID;joinReferences:ArticleID" json:"-"`
}
