package tui

import (
	"bili-kuji-management/src/db/table"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"strings"
)

func (ui *UI) exportGuidance() *tview.Grid {
	list := tview.NewList().
		AddItem("export text", "导出库存", '0', func() {
			ui.updateContent(ui.exportStock())
		}).
		AddItem("modify price", "修改价格", '1', func() {
			ui.updateContent(ui.modifyStockPricePreview())
		})

	back := tview.NewButton("Back")
	back.SetBorder(true)
	back.SetSelectedFunc(func() {
		ui.updateContent(ui.home())
	})

	return ui.contentLayout().
		AddItem(tview.NewGrid().SetBorders(true).SetRows(1, 0, 3).SetColumns(0).
			AddItem(ui.centerText("Stock & Price"), 0, 0, 1, 1, 0, 0, false).
			AddItem(list, 1, 0, 1, 1, 0, 0, true).
			AddItem(back, 2, 0, 1, 1, 0, 0, true),
			1, 1, 1, 3, 0, 0, true)
}

func (ui *UI) exportStock() *tview.Grid {
	text := ""
	stocks := ui.db.StockPreview("")

	var currentLevel int64

	for _, stock := range stocks {
		if stock.Level > currentLevel {
			currentLevel = stock.Level
			text += fmt.Sprintf("\n\n%v", ui.getLevel(stock.Level))
		}

		text += fmt.Sprintf("\n%v %v", stock.ItemName, ui.db.GetPrice(stock.ItemName).Price)
	}

	text = strings.TrimLeft(text, "\n\n")

	if err := os.WriteFile("stock.txt", []byte(text), 0644); err != nil {
		ui.logger.Error(err.Error())
	}

	content := tview.NewTextView().SetText(text).SetTextColor(tcell.ColorGreen).SetTextAlign(tview.AlignCenter)
	back := tview.NewButton("Back")
	back.SetBorder(true)
	back.SetSelectedFunc(func() {
		ui.updateContent(ui.home())
	})

	return ui.contentLayout().
		AddItem(tview.NewGrid().SetBorders(true).SetRows(1, 0, 3).SetColumns(0).
			AddItem(ui.centerText("已导出至 stock.txt").SetTextColor(tcell.ColorYellow), 0, 0, 1, 1, 0, 0, false).
			AddItem(content, 1, 0, 1, 1, 0, 0, true).
			AddItem(back, 2, 0, 1, 1, 0, 0, true),
			1, 1, 1, 3, 0, 0, true)
}

func (ui *UI) modifyStockPricePreview() *tview.Grid {
	list := tview.NewList().
		AddItem("Last", "Level: 1", '0', func() {
			ui.updateContent(ui.modifyStockPrice(1))
		}).
		AddItem("S", "Level: 2", '1', func() {
			ui.updateContent(ui.modifyStockPrice(2))
		}).
		AddItem("A", "Level: 3", '2', func() {
			ui.updateContent(ui.modifyStockPrice(3))
		}).
		AddItem("B", "Level: 4", '3', func() {
			ui.updateContent(ui.modifyStockPrice(4))
		}).
		AddItem("C", "Level: 5", '4', func() {
			ui.updateContent(ui.modifyStockPrice(5))
		}).
		AddItem("D", "Level: 6", '5', func() {
			ui.updateContent(ui.modifyStockPrice(6))
		}).
		AddItem("HIDDEN", "Level 99", '6', func() {
			ui.updateContent(ui.modifyStockPrice(99))
		})

	back := tview.NewButton("Back")
	back.SetBorder(true)
	back.SetSelectedFunc(func() {
		ui.updateContent(ui.exportGuidance())
	})

	return ui.contentLayout().
		AddItem(tview.NewGrid().SetBorders(true).SetRows(1, 0, 3).SetColumns(0).
			AddItem(ui.centerText("Select Level"), 0, 0, 1, 1, 0, 0, false).
			AddItem(list, 1, 0, 1, 1, 0, 0, true).
			AddItem(back, 2, 0, 1, 1, 0, 0, true),
			1, 1, 1, 3, 0, 0, true)

}

func (ui *UI) modifyStockPrice(level int64) *tview.Grid {
	form := tview.NewForm()
	stocks := ui.db.GetStockPreviewByLevel(level)

	for _, stock := range stocks {
		form.AddInputField(stock.ItemName, ui.db.GetPrice(stock.ItemName).Price, 10, nil, nil)
	}

	form.AddTextView("Level", ui.getLevel(level), 40, 1, true, false)
	form.AddTextView("Note", "输入你的价格后点击 Save 按钮保存", 40, 1, true, false)

	form.AddButton("Save", func() {
		for _, stock := range stocks {
			price := form.GetFormItemByLabel(stock.ItemName).(*tview.InputField).GetText()
			ui.db.UpdatePrice(&table.Price{Price: price, ItemName: stock.ItemName})
		}
		ui.updateContent(ui.modifyStockPricePreview())
	})
	return ui.contentLayout().
		AddItem(tview.NewGrid().SetBorders(true).SetRows(1, 0).SetColumns(0).
			AddItem(ui.centerText("Modify Price"), 0, 0, 1, 1, 0, 0, false).
			AddItem(form, 1, 0, 1, 1, 0, 0, true),
			1, 1, 1, 3, 0, 0, true)
}
