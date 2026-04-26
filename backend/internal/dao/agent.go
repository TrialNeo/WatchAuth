package dao

import "time"

type Agent struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Contact   string    `gorm:"type:varchar(200)" json:"contact"`
	ParentID  uint      `gorm:"not null;default:0;index" json:"parentId"`
	Level     uint8     `gorm:"not null;default:0" json:"level"` // 0=总代 1=一级 2=二级 3=三级
	Discount  float64   `gorm:"type:decimal(5,2);not null;default:100.00" json:"discount"`
	Balance   float64   `gorm:"type:decimal(12,2);not null;default:0" json:"balance"`
	Status    uint8     `gorm:"not null;default:0" json:"status"` // 0=正常 1=冻结 2=禁用
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Agent) TableName() string { return "agents" }
