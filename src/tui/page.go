package tui

import (
	"bili-kuji-management/src/db/table"
	"bili-kuji-management/src/logger"
	"fmt"
	"github.com/rivo/tview"
	"github.com/skip2/go-qrcode"
	"github.com/tidwall/gjson"
	"time"
)

func (ui *UI) home() *tview.Grid {
	list := tview.NewList().
		AddItem("transfer", "装扮转增", '0', func() {
			ui.updateContent(ui.stocksPreview())
		}).
		AddItem("account", "帐号管理", '1', func() {
			ui.updateContent(ui.accountPreview())
		}).
		AddItem("export", "导出数据", '2', func() {
			ui.updateContent(ui.exportStock())
		})

	return ui.contentLayout().
		AddItem(tview.NewGrid().SetBorders(true).SetRows(1, 0).SetColumns(0).
			AddItem(ui.centerText("Menu"), 0, 0, 1, 1, 0, 0, false).
			AddItem(list, 1, 0, 1, 1, 0, 0, true),
			1, 1, 1, 3, 0, 0, true)
}

func (ui *UI) accountPreview() *tview.Grid {
	list := tview.NewList()
	accounts := ui.db.GetAccounts()

	for k, account := range accounts {
		list.AddItem(fmt.Sprintf("%v (%v)", account.Username, account.DedeUserID), fmt.Sprintf("Online: %v", account.Online), '0'+rune(k), nil)
	}

	return ui.contentLayout().
		AddItem(tview.NewGrid().SetBorders(true).SetRows(1, 0, 3).SetColumns(0).
			AddItem(ui.centerText("Account Preview"), 0, 0, 1, 1, 0, 0, false).
			AddItem(list, 1, 0, 1, 1, 0, 0, true).
			AddItem(ui.accountGuidance(), 2, 0, 1, 1, 0, 0, true),
			1, 1, 1, 3, 0, 0, true)
}

func (ui *UI) accountGuidance() *tview.Grid {
	back := tview.NewButton("Back")
	back.SetBorder(true)
	back.SetSelectedFunc(func() {
		ui.updateContent(ui.home())
	})

	login := tview.NewButton("Login")
	login.SetBorder(true)
	login.SetSelectedFunc(func() {
		ui.updateContent(ui.qrCodeLogin())
	})

	return tview.NewGrid().SetBorders(false).SetColumns(0, 0).
		AddItem(back, 0, 0, 1, 1, 0, 0, true).
		AddItem(login, 0, 1, 1, 1, 0, 0, true)
}

func (ui *UI) stocksPreview() *tview.Grid {
	list := tview.NewList()
	stocks := ui.db.StockPreview()

	for k, stock := range stocks {
		list.AddItem(fmt.Sprintf("%v (%v)", stock.ItemName, ui.getLevel(stock.Level)), fmt.Sprintf("Stock: %v", stock.Count), '0'+rune(k), func() {
			ui.updateContent(ui.itemStockPreview(stock.ItemName))
		})
	}

	back := tview.NewButton("Back")
	back.SetBorder(true)
	back.SetSelectedFunc(func() {
		ui.updateContent(ui.home())
	})

	return ui.contentLayout().
		AddItem(tview.NewGrid().SetBorders(true).SetRows(1, 0, 3).SetColumns(0).
			AddItem(ui.centerText("Stocks Preview"), 0, 0, 1, 1, 0, 0, false).
			AddItem(list, 1, 0, 1, 1, 0, 0, true).
			AddItem(back, 2, 0, 1, 1, 0, 0, true),
			1, 1, 1, 3, 0, 0, true)
}

func (ui *UI) itemStockPreview(itemName string) *tview.Grid {
	list := tview.NewList()
	stocks := ui.db.ItemStockPreview(itemName)

	for k, stock := range stocks {
		list.AddItem(
			fmt.Sprintf("%v (%v) Time: %v",
				stock.ItemName,
				stock.ItemID,
				time.Unix(stock.GetTime, 0).Local().Format("2006-01-02 15:04:05"),
			),
			fmt.Sprintf("Owner: %v ", stock.Owner),
			'0'+rune(k),
			func() {
				ui.updateContent(ui.transferToMID(stock))
			},
		)
	}

	back := tview.NewButton("Back")
	back.SetBorder(true)
	back.SetSelectedFunc(func() {
		ui.updateContent(ui.stocksPreview())
	})

	return ui.contentLayout().
		AddItem(tview.NewGrid().SetBorders(true).SetRows(1, 0, 3).SetColumns(0).
			AddItem(ui.centerText(itemName), 0, 0, 1, 1, 0, 0, false).
			AddItem(list, 1, 0, 1, 1, 0, 0, true).
			AddItem(back, 2, 0, 1, 1, 0, 0, true),
			1, 1, 1, 3, 0, 0, true)
}

func (ui *UI) transferToMID(stock *table.Reward) *tview.Grid {
	form := tview.NewForm()

	form.AddInputField("MID", "", 20, nil, nil).
		AddButton("Save", func() {
			mid := form.GetFormItemByLabel("MID").(*tview.InputField).GetText()
			ui.updateContent(ui.transferToMIDConfirm(mid, stock))
		}).
		AddButton("Cancel", func() {
			ui.updateContent(ui.itemStockPreview(stock.ItemName))
		})

	ui.updateContent(ui.contentLayout().
		AddItem(form, 1, 1, 1, 3, 0, 0, true))

	return ui.contentLayout().
		AddItem(tview.NewGrid().SetBorders(true).SetRows(1, 0).SetColumns(0).
			AddItem(ui.centerText("Transfer To MID"), 0, 0, 1, 1, 0, 0, false).
			AddItem(form, 1, 0, 1, 1, 0, 0, true),
			1, 1, 1, 3, 0, 0, true)
}

func (ui *UI) transferToMIDConfirm(mid string, stock *table.Reward) *tview.Grid {
	resp, err := ui.request.GetCardByMid(mid)
	if err != nil {
		ui.logger.Error(err.Error())
		return ui.errorPage("TRANSFER ERROR", "Failed to get receiver info!!!")
	}

	response := gjson.Parse(resp.String())
	if response.Get("code").Int() != 0 || resp.StatusCode() != 200 {
		ui.logger.Error("get card by mid error", logger.Data("response", resp.String())...)
		return ui.errorPage("TRANSFER ERROR", "Failed to get receiver info!!!")
	}

	levelInfo := fmt.Sprintf("%v (%v/%v)",
		response.Get("card.level_info.current_level").String(),
		response.Get("card.level_info.current_exp").String(),
		response.Get("card.level_info.next_exp").String(),
	)

	form := tview.NewForm()
	form.AddTextView("MID", response.Get("card.mid").String(), 20, 1, true, false)
	form.AddTextView("Receiver", response.Get("card.name").String(), 20, 1, true, false)
	form.AddTextView("Level", levelInfo, 20, 1, true, false)
	form.AddTextView("Sign", response.Get("card.sign").String(), 20, 1, true, false)
	form.AddTextView("Item", stock.ItemName, 20, 1, true, false)
	form.AddTextView("Item ID", fmt.Sprint(stock.ItemID), 20, 1, true, false)
	form.AddTextView("Get Time", time.Unix(stock.GetTime, 0).Local().Format("2006-01-02 15:04:05"), 20, 1, true, false)
	form.AddTextView("Provider", stock.Owner, 20, 1, true, false)
	form.AddTextView("Provider MID", stock.Mid, 20, 1, true, false)

	form.AddButton("Confirm", func() {
		resp, err = ui.request.UserBoxUse(stock.ItemID, mid, ui.db.GetAccountByDedeUserID(stock.Mid))
		if err != nil {
			ui.logger.Error(err.Error())
		}

		response = gjson.Parse(resp.String())
		if response.Get("code").Int() != 0 || resp.StatusCode() != 200 {
			ui.logger.Error("transfer to mid error",
				logger.Data(
					"to_mid", mid,
					"to_user", response.Get("card.name").String(),
					"item", stock.ItemName,
					"item_id", stock.ItemID,
					"sender", stock.Owner,
					"sender_mid", stock.Mid,
					"response", resp.String(),
				)...,
			)
			ui.updateContent(ui.home())
			return
		}

		ui.logger.Info("transfer to mid success",
			logger.Data(
				"to_mid", mid,
				"to_user", response.Get("card.name").String(),
				"item", stock.ItemName,
				"item_id", stock.ItemID,
				"sender", stock.Owner,
				"sender_mid", stock.Mid,
				"response", resp.String(),
			)...,
		)

		go ui.loadStock()
		ui.updateContent(ui.home())
	})
	form.AddButton("Cancel", func() {
		ui.updateContent(ui.transferToMID(stock))
	})

	return ui.contentLayout().
		AddItem(tview.NewGrid().SetBorders(true).SetRows(1, 0).SetColumns(0).
			AddItem(ui.centerText("Transfer Confirm"), 0, 0, 1, 1, 0, 0, false).
			AddItem(form, 1, 0, 1, 1, 0, 0, true),
			1, 1, 1, 3, 0, 0, true)
}

func (ui *UI) qrCodeLogin() *tview.Grid {
	qrcodeURL := ui.login.LoginAuthCode().GetLoginUrl()

	go func() {
		ticker := time.NewTicker(time.Second * 3)

		for range ticker.C {
			if account, ok := ui.login.QrcodePoll(); ok {
				// 需要停一下不然无法验证帐号 cookies
				time.Sleep(time.Second)
				ui.db.AddNewAccount(account)
				ui.verifyAccount()
				ui.loadStock()
				ui.updateContent(ui.home())
				break
			}
		}
	}()

	qr, err := qrcode.New(qrcodeURL, qrcode.Medium)
	if err != nil {
		ui.logger.Fatal(err.Error())
	}

	textView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText(qr.ToSmallString(false))

	return ui.contentLayout().
		AddItem(textView, 1, 1, 1, 3, 0, 0, false)
}
