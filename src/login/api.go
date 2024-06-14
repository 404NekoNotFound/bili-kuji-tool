package login

import (
	"bili-kuji-management/src/db/table"
	"bili-kuji-management/src/logger"
	"net/url"
	"strconv"
	"time"
)

func (l *Login) LoginAuthCode() *Login {
	u := url.Values{}

	u.Add("appkey", "dfca71928277209b")
	u.Add("build", "1410001")
	u.Add("buvid", l.deviceID)
	u.Add("c_locale", "zh-Hans_CN")
	u.Add("channel", "bili")
	u.Add("device", "phone")
	u.Add("local_id", l.deviceID)
	u.Add("mobi_app", "android_hd")
	u.Add("platform", "android")
	u.Add("s_locale", "zh-Hans_CN")
	u.Add("spm_id", "from_spmid")
	u.Add("statistics", "{\"appId\":5,\"platform\":3,\"version\":\"1.41.0\",\"abtest\":\"\"}")
	u.Add("ts", strconv.FormatInt(time.Now().Unix(), 10))

	result := &authCodeResp{}

	r, err := l.client.R().
		SetResult(result).
		SetQueryString(signQuery(u)).
		Post("https://passport.bilibili.com/x/passport-tv-login/qrcode/auth_code")

	if err != nil {
		l.logger.Fatal(err.Error())
	}

	if result.Code != 0 || r.StatusCode() != 200 {
		l.logger.Fatal("获取登录二维码失败!", logger.Data("response", r.String())...)
	} else {
		l.logger.Info("获取登录二维码成功.")
	}

	l.authCode = result.Data.AuthCode
	l.qrcodeUrl = result.Data.Url

	return l
}

func (l *Login) QrcodePoll() (*table.Account, bool) {
	u := url.Values{}

	u.Add("appkey", "dfca71928277209b")
	u.Add("auth_code", l.authCode)
	u.Add("build", "1410001")
	u.Add("buvid", l.deviceID)
	u.Add("c_locale", "zh-Hans_CN")
	u.Add("channel", "bili")
	u.Add("device", "phone")
	u.Add("local_id", l.deviceID)
	u.Add("mobi_app", "android_hd")
	u.Add("platform", "android")
	u.Add("s_locale", "zh-Hans_CN")
	u.Add("spm_id", "from_spmid")
	u.Add("statistics", "{\"appId\":5,\"platform\":3,\"version\":\"1.41.0\",\"abtest\":\"\"}")
	u.Add("ts", strconv.FormatInt(time.Now().Unix(), 10))

	result := &qrcodePollResp{}

	r, err := l.client.R().
		SetResult(result).
		SetQueryString(signQuery(u)).
		Post("https://passport.bilibili.com/x/passport-tv-login/qrcode/poll")

	if err != nil {
		l.logger.Fatal(err.Error())
	}

	if r.StatusCode() != 200 {
		l.logger.Fatal("确认二维码状态失败!", logger.Data("response", r.String())...)
	}

	account := &table.Account{}

	if result.Code == 0 {
		l.logger.Info("登录成功.")
		for _, v := range result.Data.CookieInfo.Cookies {
			switch v.Name {
			case "SESSDATA":
				account.SESSDATA = v.Value
			case "bili_jct":
				account.BiliJct = v.Value
			case "DedeUserID":
				account.DedeUserID = v.Value
			case "DedeUserID__ckMd5":
				account.DedeUserIDCkMd5 = v.Value
			case "sid":
				account.Sid = v.Value
			}
		}

		account.Buvid = l.deviceID
		account.AccessKey = result.Data.TokenInfo.AccessToken
		account.RefreshToken = result.Data.TokenInfo.RefreshToken

		return account, true
	}

	l.logger.Info("二维码当前状态", logger.Data("response", r.String())...)

	return nil, false
}

func (l *Login) LoginExitV2(account *table.Account) *Login {
	u := url.Values{}

	u.Add("biliCSRF", account.BiliJct)

	result := &loginExitV2Resp{}

	_, err := l.client.R().
		SetResult(result).
		SetQueryString(u.Encode()).
		SetCookies(l.GetCookies(account)).
		Post("https://passport.bilibili.com/login/exit/v2")

	if err != nil {
		l.logger.Error(err.Error())
	}

	return l
}
