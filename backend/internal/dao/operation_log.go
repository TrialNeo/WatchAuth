package dao

import "time"

type OperationLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AdminID   uint      `gorm:"index;not null;default:0" json:"adminId"`
	AdminName string    `gorm:"type:varchar(50)" json:"adminName"`
	Action    string    `gorm:"type:varchar(50);index;not null" json:"action"`
	Target    string    `gorm:"type:varchar(100)" json:"target"`
	TargetID  string    `gorm:"type:varchar(36)" json:"targetId"`
	Detail    string    `gorm:"type:text" json:"detail"`
	IP        string    `gorm:"type:varchar(45)" json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
}

func (OperationLog) TableName() string { return "operation_logs" }
