package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func newApp() *tview.Application {
	return tview.NewApplication()
}

func defaultHeaderLeft() *tview.TextView {
	return tview.NewTextView().SetText("\n帐号数: loading...").SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorRed)
}

func defaultHeaderMiddleLeft() *tview.TextView {
	return tview.NewTextView().SetText("\n库存: loading...").SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorRed)
}

func defaultHeaderMiddleRight() *tview.TextView {
	return tview.NewTextView().SetText(fmt.Sprint("\n奖池: ", activityName)).SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorLightPink)
}

func defaultHeaderRight() *tview.TextView {
	return tview.NewTextView().SetText("\nVersion: " + version).SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorBlue)
}

func defaultFooterLeft() *tview.TextView {
	return tview.NewTextView().SetText("Sponsor: 装扮小姐姐").SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorPurple)
}

func defaultFooterMiddle() *tview.TextView {
	return tview.NewTextView().SetText("Author: Mika_喵").SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorDarkGray)
}

func defaultFooterRight() *tview.TextView {
	return tview.NewTextView().SetText("https://space.bilibili.com/1972716772").SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorYellow)
}

func (ui *UI) newHeader() *tview.Grid {
	return tview.NewGrid().SetRows(0).SetColumns(0, 0, 0, 0).SetBorders(false).
		AddItem(ui.headerLeft, 0, 0, 1, 1, 0, 0, false).
		AddItem(ui.headerMiddleLeft, 0, 1, 1, 1, 0, 0, false).
		AddItem(ui.headerMiddleRight, 0, 2, 1, 1, 0, 0, false).
		AddItem(ui.headerRight, 0, 3, 1, 1, 0, 0, false)
}

func (ui *UI) newFooter() *tview.Grid {
	return tview.NewGrid().SetRows(0).SetColumns(0, 0, 0).SetBorders(false).
		AddItem(ui.footerLeft, 0, 0, 1, 1, 0, 0, false).
		AddItem(ui.footerMiddle, 0, 1, 1, 1, 0, 0, false).
		AddItem(ui.footerRight, 0, 2, 1, 1, 0, 0, false)
}

func (ui *UI) newLayout() *tview.Grid {
	return tview.NewGrid().SetRows(3, 0, 1).SetBorders(true).
		AddItem(ui.header, 0, 0, 1, 1, 0, 0, false).
		AddItem(ui.content, 1, 0, 1, 1, 0, 0, false).
		AddItem(ui.footer, 2, 0, 1, 1, 0, 0, false)
}

func (ui *UI) contentLayout() *tview.Grid {
	return tview.NewGrid().SetBorders(false).
		SetRows(2, 0).
		SetColumns(0, 0, 0, 0, 0).
		AddItem(tview.NewTextView(), 0, 0, 1, 5, 0, 0, false).
		AddItem(tview.NewTextView(), 1, 0, 1, 1, 0, 0, false).
		AddItem(tview.NewTextView(), 1, 4, 1, 1, 0, 0, false)
}

func (ui *UI) centerText(title string) *tview.TextView {
	return tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(title)
}

// 更新 Content 中的内容
func (ui *UI) updateContent(newContent *tview.Grid) *UI {
	ui.layout.RemoveItem(ui.content)
	ui.content = newContent
	ui.app.SetFocus(ui.content)
	ui.layout.AddItem(ui.content, 1, 0, 1, 1, 0, 0, false)
	ui.app.ForceDraw()

	return ui
}

func (ui *UI) updateAccountCount() *UI {
	ui.app.QueueUpdateDraw(func() {
		count := len(ui.db.GetAccountsOnline())
		if count == 0 {
			ui.headerLeft.SetText("\nAccount: 0").SetTextColor(tcell.ColorRed)
		} else {
			ui.headerLeft.SetText(fmt.Sprintf("\nAccount: %v", count)).SetTextColor(tcell.ColorGreen)
		}
	})
	return ui
}

func (ui *UI) updateStockCount() *UI {
	ui.app.QueueUpdateDraw(func() {
		count := len(ui.db.GetStocks())
		if count == 0 {
			ui.headerMiddleLeft.SetText("\n库存: 0").SetTextColor(tcell.ColorRed)
		} else {
			ui.headerMiddleLeft.SetText(fmt.Sprintf("\n库存: %v", count)).SetTextColor(tcell.ColorGreen)
		}
	})
	return ui
}

func (ui *UI) errorPage(heading, text string) *tview.Grid {
	content := tview.NewTextView().SetText(text).SetTextColor(tcell.ColorRed).SetTextAlign(tview.AlignCenter)
	back := tview.NewButton("Back")
	back.SetBorder(true)
	back.SetSelectedFunc(func() {
		ui.updateContent(ui.home())
	})

	return ui.contentLayout().
		AddItem(tview.NewGrid().SetBorders(true).SetRows(1, 0, 3).SetColumns(0).
			AddItem(ui.centerText(heading), 0, 0, 1, 1, 0, 0, false).
			AddItem(content, 1, 0, 1, 1, 0, 0, true).
			AddItem(back, 2, 0, 1, 1, 0, 0, true),
			1, 1, 1, 3, 0, 0, true)
}
