package dao

import (
	"time"
)

type Version struct {
	Id          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Appid       string    `gorm:"type:varchar(36);not null" json:"appid"`
	Version     string    `gorm:"not null" json:"version"`
	Desc        string    `json:"desc"`
	Sign        string    `gorm:"not null" json:"sign"` //使用MD5
	ForceUpdate bool      `gorm:"not null" json:"force_update"`
	Status      bool      `gorm:"not null" json:"status"`
	PatchUrl    string    `gorm:"not null" json:"patch_url"`
	CreatedAt   time.Time `json:"updatedTime"`         //这个其实是更新时间
	DeletedAt   time.Time `json:"deletedAt,omitempty"` // 软删除字段
}
