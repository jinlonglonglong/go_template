package models

type User struct {
	ID      int64  `json:"id" gorm:"column:id"`           // 主键id
	Address string `json:"address" gorm:"column:address"` // 用户地址
	Balance string `json:"balance" gorm:"column:balance"` // 用户代币余额
}

func (m *User) TableName() string {
	return "user"
}
