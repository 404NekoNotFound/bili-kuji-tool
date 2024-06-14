package request

import (
	"github.com/go-resty/resty/v2"
	"net/url"
	"strconv"
)

func (r *Request) AssistList(account *Account) (*resty.Response, error) {
	u := url.Values{}
	u.Add("activity_id", strconv.Itoa(r.activityID))
	return r.client.R().
		SetQueryString(u.Encode()).
		SetCookies(account.Cookies).
		Get("https://api.bilibili.com/x/garb/kuji/v5/assist/list")
}

func (r *Request) UserBoxList(state int, account *Account) (*resty.Response, error) {
	u := url.Values{}
	u.Add("kuji_act_id", strconv.Itoa(r.activityID))
	u.Add("state", strconv.Itoa(state))

	return r.client.R().
		SetQueryString(u.Encode()).
		SetCookies(account.Cookies).
		Get("https://api.bilibili.com/x/garb/kuji/v5/user/box/list")
}

func (r *Request) UserBoxUse(id int64, toMid string, account *Account) (*resty.Response, error) {
	u := url.Values{}
	u.Add("id", strconv.FormatInt(id, 10))
	u.Add("to_mid", toMid)
	u.Add("csrf", account.BiliJct)

	return r.client.R().
		SetQueryString(u.Encode()).
		SetCookies(account.Cookies).
		Post("https://api.bilibili.com/x/garb/kuji/v5/user/box/use")
}
