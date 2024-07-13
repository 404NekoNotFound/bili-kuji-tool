package db

import (
	"bili-kuji-management/src/db/table"
	"bili-kuji-management/src/request"
	"net/http"
	"sort"
)

func (db *DB) GetAccounts() []*request.Account {
	var accounts []table.Account

	if err := db.db.Find(&accounts).Error; err != nil {
		db.log.Error(err.Error())
	}

	return db.accountFormat(accounts)
}

func (db *DB) GetAccountsOnline() []*request.Account {
	var accounts []table.Account

	if err := db.db.Where(&table.Account{Online: true}).Find(&accounts).Error; err != nil {
		db.log.Error(err.Error())
	}

	return db.accountFormat(accounts)
}

func (db *DB) GetStocks() []*table.Reward {
	var rewards []*table.Reward

	if err := db.db.Model(&table.Reward{}).Find(&rewards).Error; err != nil {
		db.log.Error(err.Error())
	}

	return rewards
}

func (db *DB) StockPreview(keyword string) []*table.RewardPreview {
	var result []*table.RewardPreview

	if err := db.db.Model(&table.Reward{}).
		Select("item_name", "level", "COUNT(level) AS count").
		Where("item_name LIKE ?", "%"+keyword+"%").
		Group("item_name").
		Order("level").
		Find(&result).Error; err != nil {
		db.log.Error(err.Error())
	}

	return result
}

func (db *DB) ItemStockPreview(itemName string) []*table.Reward {
	var reward []*table.Reward

	if err := db.db.Where(&table.Reward{ItemName: itemName}).Order("get_time").Find(&reward).Error; err != nil {
		db.log.Error(err.Error())
	}

	return reward
}

func (db *DB) GetAccountByDedeUserID(DedeUserID string) *request.Account {
	var account table.Account

	if err := db.db.Where(&table.Account{DedeUserID: DedeUserID}).Find(&account).Error; err != nil {
		db.log.Error(err.Error())
	}

	return db.accountFormat([]table.Account{account})[0]
}

func (db *DB) accountFormat(accounts []table.Account) []*request.Account {
	var users []*request.Account

	for _, account := range accounts {
		users = append(users, &request.Account{
			AccessKey:       account.AccessKey,
			RefreshToken:    account.RefreshToken,
			SESSDATA:        account.SESSDATA,
			BiliJct:         account.BiliJct,
			DedeUserID:      account.DedeUserID,
			DedeUserIDCkMd5: account.DedeUserIDCkMd5,
			Sid:             account.Sid,
			Buvid:           account.Buvid,
			Username:        account.Username,
			Online:          account.Online,
			Cookies: []*http.Cookie{
				{Name: "SESSDATA", Value: account.SESSDATA},
				{Name: "bili_jct", Value: account.BiliJct},
				{Name: "DedeUserID", Value: account.DedeUserID},
				{Name: "DedeUserID__ckMd5", Value: account.DedeUserIDCkMd5},
				{Name: "sid", Value: account.Sid},
				{Name: "Buvid", Value: account.Buvid},
			},
		})
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].DedeUserID < users[j].DedeUserID
	})

	return users
}

func (db *DB) GetStockPreviewByLevel(level int64) []*table.RewardPreview {
	var result []*table.RewardPreview

	if err := db.db.Model(&table.Reward{}).
		Select("item_name", "level", "COUNT(level) AS count").
		Where(&table.Reward{Level: level}).
		Group("item_name").
		Order("level").
		Find(&result).Error; err != nil {
		db.log.Error(err.Error())
	}

	return result
}

func (db *DB) GetPrice(itemName string) *table.Price {
	prices := &table.Price{}

	if err := db.db.Where(&table.Price{ItemName: itemName}).Find(prices).Error; err != nil {
		db.log.Error(err.Error())
	}

	return prices
}
