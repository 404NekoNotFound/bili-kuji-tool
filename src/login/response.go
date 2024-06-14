package login

type authCodeResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		Url      string `json:"url"`
		AuthCode string `json:"auth_code"`
	} `json:"data"`
}

type qrcodePollResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		IsNew        bool   `json:"is_new"`
		Mid          int    `json:"mid"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenInfo    struct {
			Mid          int    `json:"mid"`
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    int    `json:"expires_in"`
		} `json:"token_info"`
		CookieInfo struct {
			Cookies []struct {
				Name     string `json:"name"`
				Value    string `json:"value"`
				HttpOnly int    `json:"http_only"`
				Expires  int    `json:"expires"`
			} `json:"cookies"`
			Domains []string `json:"domains"`
		} `json:"cookie_info"`
		Sso []string `json:"sso"`
	} `json:"data"`
}

type loginExitV2Resp struct {
	Code   int  `json:"code"`
	Status bool `json:"status"`
	Ts     int  `json:"ts"`
	Data   struct {
		RedirectUrl string `json:"redirectUrl"`
	} `json:"data"`
}
