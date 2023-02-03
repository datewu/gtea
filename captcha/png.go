package captcha

import (
	"bytes"
	"image"
	"image/png"
	"math"
	"math/rand"
	"sync"
)

func randomN(m int) int {
	if m < 2 {
		return rand.Intn(10)
	}
	m = m - 1
	r := rand.Int31n(9 * int32(math.Pow10(m)))
	return int(r) + int(math.Pow10(m))
}

func genPNG(number int) []byte {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{captchaWidth + 2*captchaXPadding +
		(totalCaptchaDs-1)*captchaGap, captchaHeight}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	ds := digits(number, totalCaptchaDs)
	var wg sync.WaitGroup
	for i, d := range ds {
		wg.Add(1)
		go func(n, offset int) {
			defer wg.Done()
			drawMahjong(n, offset, img)
		}(d, captchaXPadding+i*(singleCaptchaWidth+captchaGap))
	}
	wg.Wait()
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func digits(number, hight int) []int {
	if hight < 2 {
		return []int{number}
	}
	var digits []int
	for i := hight - 1; i >= 0; i-- {
		d := number / int(math.Pow10(i))
		number -= d * int(math.Pow10(i))
		digits = append(digits, d)
	}
	return digits
}
