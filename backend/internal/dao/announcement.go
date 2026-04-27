package dao

import "time"

type Announcement struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"type:varchar(200);not null" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Status    uint8     `gorm:"not null;default:0" json:"status"` // 0=发布 1=草稿
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Announcement) TableName() string { return "announcements" }
