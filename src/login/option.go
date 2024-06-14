package login

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Login struct {
	client         *resty.Client
	logger         *zap.Logger
	authCode       string
	qrcodeUrl      string
	androidVersion string
	deviceName     string
	deviceID       string
}

type Option func(l *Login)

func (opt Option) Apply(l *Login) {
	opt(l)
}

func Ready(opts ...Option) *Login {
	login := &Login{
		client:   resty.New().SetHeader("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		deviceID: buvid(),
	}

	for _, o := range opts {
		o.Apply(login)
	}

	// 填充设备版本
	login.androidVersion = "13.0"
	login.deviceName = "Pixel 7"
	//login.logger.Warn("已设置设备ID", logger.Data("buvid", login.deviceID)...)

	return login
}

func WithLogger(logger *zap.Logger) Option {
	return func(l *Login) {
		l.logger = logger
	}
}
