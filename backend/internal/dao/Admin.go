package dao

import (
	"time"
)

type Admin struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // json:"-" 表示返回给前端时隐藏密码
	LastIp    string    `gorm:"type:varchar(255);" json:"lastIp"`
	Ips       string    `gorm:"type:text" json:"ips"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type App struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AppID      string    `gorm:"type:varchar(36);uniqueIndex;not null" json:"appId"`
	AppName    string    `gorm:"type:text;uniqueIndex';not null" json:"appName"`
	Desc       string    `gorm:"type:text;not null" json:"desc"`
	Version    string    `gorm:"type:text;not null" json:"version"`
	FeeType    uint8     `gorm:"not null" json:"feeType"`
	Status     uint8     `gorm:"not null" json:"status"`
	Fee        float64   `gorm:"not null" json:"fee"`
	EncType    uint8     `gorm:"not null" json:"encType"`
	SecretKeys string    `gorm:"type:text;" json:"secret"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	DeletedAt  time.Time `json:"deletedAt,omitempty"` // GORM 软删除字段
}
