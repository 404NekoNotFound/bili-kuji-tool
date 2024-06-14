package request

import "net/http"

type Account struct {
	AccessKey       string
	RefreshToken    string
	SESSDATA        string
	BiliJct         string
	DedeUserID      string
	DedeUserIDCkMd5 string
	Sid             string
	Buvid           string
	Username        string
	Online          bool
	Cookies         []*http.Cookie
}

func (r *Request) GetActivityID() int {
	return r.activityID
}
