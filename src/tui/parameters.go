package tui

import (
	"bili-kuji-management/src/db/table"
	"bili-kuji-management/src/logger"
	"bili-kuji-management/src/request"
	"github.com/tidwall/gjson"
	"sync"
)

func (ui *UI) verifyAccount() *UI {
	accounts := ui.db.GetAccounts()
	wg := &sync.WaitGroup{}
	wg.Add(len(accounts))

	for _, account := range accounts {
		go ui.verifyAccountRequest(account, wg)
	}

	wg.Wait()

	ui.updateAccountCount()
	return ui
}

func (ui *UI) verifyAccountRequest(account *request.Account, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := ui.request.Nav(account)
	if err != nil {
		ui.logger.Error(err.Error())
		return
	}

	response := gjson.Parse(resp.String())

	if resp.StatusCode() != 200 || response.Get("code").Int() != 0 {
		ui.logger.Error("failed to verify account",
			logger.Data(
				"mid", account.DedeUserID,
				"username", account.Username,
				"response", response.String(),
			)...,
		)
		return
	}

	ui.db.UpdatesAccount(&table.Account{
		DedeUserID: account.DedeUserID,
		Username:   response.Get("data.uname").String(),
		Online:     true,
	})

	ui.logger.Info("account is online",
		logger.Data(
			"mid", account.DedeUserID,
			"username", account.Username,
		)...,
	)
}

func (ui *UI) loadStock() *UI {
	ui.db.DefaultRewards()
	accounts := ui.db.GetAccountsOnline()
	wg := sync.WaitGroup{}
	wg.Add(len(accounts))

	for _, account := range accounts {
		go ui.loadStockRequest(account, &wg)
	}

	wg.Wait()
	ui.updateStockCount()
	return ui
}

func (ui *UI) loadStockRequest(account *request.Account, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := ui.request.UserBoxList(0, account)
	if err != nil {
		ui.logger.Error(err.Error())
		return
	}

	response := gjson.Parse(resp.String())

	if resp.StatusCode() != 200 || response.Get("code").Int() != 0 {
		ui.logger.Error("failed to load stock",
			logger.Data(
				"mid", account.DedeUserID,
				"username", account.Username,
				"response", response.String(),
			)...,
		)
		return
	}

	var rewards []*table.Reward

	levelLists := response.Get("data.list").Array()

	for _, levelList := range levelLists {
		level := levelList.Get("level").Int()
		list := levelList.Get("list").Array()

		for _, item := range list {
			name := item.Get("name").String()
			id := item.Get("id").Int()
			time := item.Get("get_time").Int()
			source := item.Get("source").String()

			rewards = append(rewards, &table.Reward{
				Level:    level,
				ItemID:   id,
				ItemName: name,
				GetTime:  time,
				Source:   source,
				Owner:    account.Username,
				Mid:      account.DedeUserID,
			})
		}
	}

	ui.db.InsertRewards(rewards)

	ui.logger.Info("stock loaded",
		logger.Data(
			"mid", account.DedeUserID,
			"username", account.Username,
			"count", len(rewards),
		)...,
	)
}

func (ui *UI) getLevel(level int64) string {
	return levels[level]
}

var (
	levels = map[int64]string{
		1:  "Last",
		2:  "S",
		3:  "A",
		4:  "B",
		5:  "C",
		6:  "D",
		7:  "E",
		8:  "F",
		9:  "H",
		10: "G",
		99: "HIDDEN",
	}
)
