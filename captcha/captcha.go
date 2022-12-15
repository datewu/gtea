package captcha

import (
	"encoding/base64"
	"os"
	"strconv"
	"time"

	"github.com/datewu/gtea/jsonlog"
	"github.com/datewu/security"
)

// Captcha ...
type Captcha struct {
	Tag int    `json:"captcha"`
	MD5 string `json:"check"`
	TS  int64  `json:"timestamp"` // now.Unix() second
}

func md5Captcha(tag int, ts int64) string {
	t := strconv.Itoa(tag + int(ts))
	return security.ToHexString(
		security.Md5WithTag(t, []byte(os.Getenv("CAPTCHA_SECRET"))))
}

func (c Captcha) OK() bool {
	interval := time.Now().Unix() - c.TS
	if int(interval) > 20 {
		return false
	}
	jsonlog.Info("captcha interval between client and server",
		map[string]interface{}{"interval": interval, "captcha": c.Tag})
	return c.MD5 == md5Captcha(c.Tag, c.TS)
}

func (c Captcha) PNG() string {
	b64 := "data:image/png;base64,"
	data := base64.StdEncoding.EncodeToString(genPNG(c.Tag))
	return b64 + data
}

// NewCaptcha ...
func NewCaptcha() Captcha {
	n := randomN(totalCaptchaDs)
	ts := time.Now().Unix()
	return Captcha{
		Tag: n,
		TS:  ts,
		MD5: md5Captcha(n, ts),
	}
}
