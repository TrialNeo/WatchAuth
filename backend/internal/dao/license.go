package dao

import "time"

type License struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	LicenseKey string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"licenseKey"`
	AppID      string    `gorm:"type:varchar(36);not null;index" json:"appId"`
	MachineID  int       `gorm:"index;default:0" json:"machineId"`
	UserID     uint      `gorm:"index;default:0" json:"userId"`
	ExpireAt   time.Time `json:"expireAt"`
	IssuedAt   time.Time `json:"issuedAt"`
	Status     uint8     `gorm:"not null;default:0" json:"status"` // 0=有效 1=已吊销 2=已过期
}

func (License) TableName() string { return "licenses" }
