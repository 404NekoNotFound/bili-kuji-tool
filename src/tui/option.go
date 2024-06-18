package tui

import (
	"bili-kuji-management/src/db"
	"bili-kuji-management/src/login"
	"bili-kuji-management/src/request"
	"github.com/rivo/tview"
	"go.uber.org/zap"
)

const (
	activityName = "星时遗迹"
	activityId   = 32
	version      = "v1.1.1"
)

type UI struct {
	app *tview.Application

	// Module
	layout  *tview.Grid
	header  *tview.Grid
	footer  *tview.Grid
	content *tview.Grid

	// header
	headerLeft        *tview.TextView // 帐号数
	headerMiddleLeft  *tview.TextView // 库存数
	headerMiddleRight *tview.TextView // 活动名称
	headerRight       *tview.TextView // 版本

	// footer
	footerLeft   *tview.TextView // 赞助
	footerMiddle *tview.TextView // 作者
	footerRight  *tview.TextView // 链接

	db      *db.DB
	login   *login.Login
	logger  *zap.Logger
	request *request.Request
}

type Option func(ui *UI)

func (opt Option) Apply(ui *UI) {
	opt(ui)
}

func Ready(opts ...Option) *UI {
	ui := &UI{
		app:               newApp(),
		headerLeft:        defaultHeaderLeft(),
		headerMiddleLeft:  defaultHeaderMiddleLeft(),
		headerMiddleRight: defaultHeaderMiddleRight(),
		headerRight:       defaultHeaderRight(),
		footerLeft:        defaultFooterLeft(),
		footerMiddle:      defaultFooterMiddle(),
		footerRight:       defaultFooterRight(),
	}

	for _, o := range opts {
		o.Apply(ui)
	}

	ui.header = ui.newHeader()
	ui.footer = ui.newFooter()
	ui.content = ui.home()
	ui.layout = ui.newLayout()

	ui.db = db.Connect(ui.logger)

	ui.login = login.Ready(
		login.WithLogger(ui.logger),
	)

	ui.request = request.New(
		request.WithActivityID(activityId),
	)

	ui.db.DefaultAccountOffline()
	ui.initParams()

	return ui
}

func WithLogger(logger *zap.Logger) Option {
	return func(ui *UI) {
		ui.logger = logger
	}
}

func (ui *UI) Run() {
	if err := ui.app.SetRoot(ui.layout, true).SetFocus(ui.content).EnableMouse(true).Run(); err != nil {
		ui.logger.Fatal(err.Error())
	}
}

func (ui *UI) initParams() {
	go func() {
		ui.verifyAccount()
		ui.loadStock()
	}()
}
