package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"strings"
)

func (ui *UI) exportStock() *tview.Grid {
	text := ""
	stocks := ui.db.StockPreview()

	var currentLevel int64

	for _, stock := range stocks {
		if stock.Level > currentLevel {
			currentLevel = stock.Level
			text += fmt.Sprintf("\n\n%v", ui.getLevel(stock.Level))
		}

		text += fmt.Sprintf("\n%v", stock.ItemName)
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
