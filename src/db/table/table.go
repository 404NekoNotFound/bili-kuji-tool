package table

import "gorm.io/gorm"

// Account 帐号
type Account struct {
	gorm.Model

	AccessKey       string
	RefreshToken    string
	SESSDATA        string
	BiliJct         string
	DedeUserID      string `gorm:"primary_key"`
	DedeUserIDCkMd5 string
	Sid             string
	Buvid           string
	Username        string
	Online          bool
}

// Reward 库存
type Reward struct {
	gorm.Model

	Level    int64
	ItemID   int64
	ItemName string `gorm:"primary_key"`
	GetTime  int64
	Source   string
	Owner    string
	Mid      string
}

// Price 自定义价格
type Price struct {
	gorm.Model

	ItemName string `gorm:"primary_key"`
	Price    string
}

// RewardPreview 库存简介
type RewardPreview struct {
	ItemName string
	Count    int64
	Level    int64
}
