package request

import (
	"github.com/go-resty/resty/v2"
	"net/url"
)

func (r *Request) Nav(account *Account) (*resty.Response, error) {
	return r.client.R().
		SetCookies(account.Cookies).
		Get("https://api.bilibili.com/x/web-interface/nav")
}

func (r *Request) GetCardByMid(mid string) (*resty.Response, error) {
	u := url.Values{}
	u.Add("mid", mid)

	return r.client.R().
		ForceContentType("application/json").
		SetQueryString(u.Encode()).
		Get("https://account.bilibili.com/api/member/getCardByMid")
}
