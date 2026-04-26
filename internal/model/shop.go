package model

type Shop struct {
	ID       int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string  `gorm:"type:varchar(128)" json:"name"`
	TypeID   int64   `gorm:"index" json:"typeId"`
	Images   string  `gorm:"type:varchar(1024)" json:"images"`
	Area     string  `gorm:"type:varchar(128)" json:"area"`
	Address  string  `gorm:"type:varchar(255)" json:"address"`
	Score    float64 `gorm:"type:decimal(2,1)" json:"score"` // 评分
	AvgPrice int64   `gorm:"type:int" json:"avgPrice"`       // 人均价格(存分，不存元)
}
