package dao

type SystemConfig struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Key   string `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Value string `gorm:"type:text;not null" json:"value"`
}

func (SystemConfig) TableName() string { return "system_configs" }
