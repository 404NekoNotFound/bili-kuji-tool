package db

import (
	"bili-kuji-management/src/db/table"
)

func (db *DB) DefaultRewards() *DB {
	if err := db.db.Migrator().DropTable(&table.Reward{}); err != nil {
		db.log.Error(err.Error())
	}

	if err := db.db.AutoMigrate(&table.Reward{}); err != nil {
		db.log.Error(err.Error())
	}

	return db
}

func (db *DB) DefaultAccountOffline() *DB {
	if err := db.db.Where(table.Account{Online: true}).Model(&table.Account{}).UpdateColumn("online", false).Error; err != nil {
		db.log.Error(err.Error())
	}

	return db
}

func (db *DB) AddNewAccount(account *table.Account) *DB {
	if err := db.db.
		Where(&table.Account{DedeUserID: account.DedeUserID}).
		Assign(&table.Account{
			AccessKey:       account.AccessKey,
			RefreshToken:    account.RefreshToken,
			SESSDATA:        account.SESSDATA,
			BiliJct:         account.BiliJct,
			DedeUserID:      account.DedeUserID,
			DedeUserIDCkMd5: account.DedeUserIDCkMd5,
			Sid:             account.Sid,
			Buvid:           account.Buvid,
		}).FirstOrCreate(account).Error; err != nil {
		db.log.Error(err.Error())
	}

	return db
}

func (db *DB) UpdatesAccount(account *table.Account) {
	if err := db.db.
		Model(&table.Account{}).
		Where(&table.Account{
			DedeUserID: account.DedeUserID,
		}).
		Updates(account).Error; err != nil {
		db.log.Error(err.Error())
	}
}

func (db *DB) InsertRewards(rewards []*table.Reward) {
	if err := db.db.CreateInBatches(&rewards, 100).Error; err != nil {
		db.log.Error(err.Error())
	}
}
