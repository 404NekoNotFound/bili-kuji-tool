package login

import (
	"bili-kuji-management/src/db/table"
	"net/http"
)

func (l *Login) GetLoginUrl() string {
	return l.qrcodeUrl
}

func (l *Login) GetAppKey() string {
	return "dfca71928277209b"
}

func (l *Login) GetAppSalt() string {
	return "b5475a8825547a4fc26c7d518eaaa02e"
}

func (l *Login) GetCookies(account *table.Account) []*http.Cookie {
	return []*http.Cookie{
		{Name: "SESSDATA", Value: account.SESSDATA},
		{Name: "bili_jct", Value: account.BiliJct},
		{Name: "DedeUserID", Value: account.DedeUserID},
		{Name: "DedeUserID__ckMd5", Value: account.DedeUserIDCkMd5},
		{Name: "sid", Value: account.Sid},
	}
}
