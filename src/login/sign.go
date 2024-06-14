package login

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
)

func signQuery(u url.Values) string {
	rawQuery := u.Encode()
	sign := strMd5(fmt.Sprintf("%v%v", rawQuery, "b5475a8825547a4fc26c7d518eaaa02e"))
	params := fmt.Sprintf("%v&sign=%v", rawQuery, sign)

	return params
}

func strMd5(str string) (retMd5 string) {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
