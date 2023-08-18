package captcha

import (
	"encoding/base64"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/datewu/gtea/jsonlog"
	"github.com/datewu/security"
)

const (
	captchaWidth       = 240
	captchaXPadding    = 5
	captchaGap         = 10
	totalCaptchaDs     = 6
	singleCaptchaWidth = captchaWidth / totalCaptchaDs
	captchaHeight      = 80
)

func randomN(m int) int {
	if m < 2 {
		return rand.Intn(10)
	}
	m = m - 1
	r := rand.Int31n(9 * int32(math.Pow10(m)))
	return int(r) + int(math.Pow10(m))
}

type Pnger interface {
	PNG(int) []byte
}

// Captcha ...
type Captcha struct {
	Tag int    `json:"captcha"`
	MD5 string `json:"check"`
	TS  int64  `json:"timestamp"` // now.Unix() second
	Pic Pnger  `json:"-"`
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
	data := base64.StdEncoding.EncodeToString(c.Pic.PNG(c.Tag))
	return b64 + data
}

// NewCaptcha ...
func NewCaptcha(p Pnger) Captcha {
	n := randomN(totalCaptchaDs)
	ts := time.Now().Unix()
	c := Captcha{
		Tag: n,
		TS:  ts,
		MD5: md5Captcha(n, ts),
	}
	if p == nil {
		c.Pic = MajongPng{}
	}
	return c
}
