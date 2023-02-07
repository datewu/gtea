package captcha

import (
	"image"
	"image/color"
	"math/rand"
)

const (
	captchaWidth       = 240
	captchaXPadding    = 5
	captchaGap         = 10
	totalCaptchaDs     = 6
	singleCaptchaWidth = captchaWidth / totalCaptchaDs
	captchaHeight      = 80
)

var (
	colorOuterBlue   = color.RGBA{55, 55, 180, 255}
	colorMiddleGreen = color.RGBA{40, 90, 70, 255}
	colorCenterRed   = color.RGBA{190, 30, 75, 255}
)

func squareDistance(m, n, x, y int) int {
	return (m-x)*(m-x) + (n-y)*(n-y)
}

func drawNine(x, y int, img *image.RGBA) {
	img.Set(x-1, y-1, colorOuterBlue)
	img.Set(x, y-1, colorOuterBlue)
	img.Set(x+1, y-1, colorOuterBlue)

	img.Set(x-1, y, colorOuterBlue)
	img.Set(x, y, image.Black)
	img.Set(x+1, y, colorOuterBlue)

	img.Set(x-1, y+1, colorOuterBlue)
	img.Set(x, y+1, colorOuterBlue)
	img.Set(x+1, y+1, colorOuterBlue)
}

func drawNormalRingAndCircle(x, y, r int, color color.Color, img *image.RGBA) {
	thick := 2 * r
	drawRing(x, y, r+3*thick, thick, color, img)
	drawRing(x, y, r+thick, thick, color, img)
	drawCircle(x, y, r, color, img)
}

func drawRing(x, y, rIn, thick int, color color.Color, img *image.RGBA) {
	for m := x - rIn - thick; m < x+rIn+thick; m++ {
		for n := y - rIn - thick; n < y+rIn+thick; n++ {
			if m > x-rIn+3 && m < x+rIn-3 && n > y-rIn+3 && n < y+rIn-3 {
				continue
			}
			dis := squareDistance(m, n, x, y)
			if dis > rIn*rIn && dis < (rIn+thick)*(rIn+thick) {
				img.Set(m, n, color)
			}

		}
	}
}

func drawCircle(x, y, r int, color color.Color, img *image.RGBA) {
	for m := x - r; m < x+r; m++ {
		for n := y - r; n < y+r; n++ {
			if squareDistance(m, n, x, y) < r*r {
				img.Set(m, n, color)
			}
		}
	}
}

func drawMahjong(n, xoffset int, img *image.RGBA) {
	for i := 0; i < 5; i++ {
		randx := rand.Intn(singleCaptchaWidth)
		randy := rand.Intn(captchaHeight)
		drawNine(randx+xoffset, randy, img)
	}
	r := 1
	switch n {
	case 0:
		x := singleCaptchaWidth/2 + xoffset
		y := captchaHeight / 2
		thick := 4
		width := singleCaptchaWidth / 5 * 4
		height := captchaHeight / 7 * 4
		left := x - width/2 - thick
		right := x + width/2 + thick
		up := y - height/2 - thick
		down := y + height/2 + thick
		for m := left; m < right; m++ {
			for n := up; n < up+thick; n++ {
				if n%2 == 1 {
					if m > left+2 && m < right-2 {
						continue
					}
				}
				img.Set(m, n, colorOuterBlue)
			}
			for n := down; n > down-thick; n-- {
				if n%2 == 1 {
					if m > left+2 && m < right-2 {
						continue
					}
				}
				img.Set(m, n, colorOuterBlue)
			}
		}
		for n := up; n < down; n++ {
			for m := left; m < left+thick; m++ {
				if m%2 == 1 {
					if n > up+4 && n < down-4 {
						continue
					}
				}
				img.Set(m, n, colorOuterBlue)
			}
			for m := right; m > right-thick; m-- {
				if m%2 == 1 {
					if n > up+4 && n < down-4 {
						continue
					}
				}
				img.Set(m, n, colorOuterBlue)
			}
		}
	case 1:
		x := singleCaptchaWidth/2 + xoffset
		y := captchaHeight / 2
		thick := 4 * r
		drawRing(x, y, 16, thick, colorOuterBlue, img)
		drawRing(x, y, 8, thick, colorMiddleGreen, img)
		drawCircle(x, y, 2*r, colorCenterRed, img)
	case 2:
		x := singleCaptchaWidth/2 + xoffset
		y1 := captchaHeight/2 - captchaHeight/6
		y2 := captchaHeight/2 + captchaHeight/6
		thick := 2 * r
		drawRing(x, y1, r+5*thick, thick, colorOuterBlue, img)
		drawNormalRingAndCircle(x, y1, r, colorOuterBlue, img)

		drawRing(x, y2, r+5*thick, thick, colorMiddleGreen, img)
		drawNormalRingAndCircle(x, y2, r, colorMiddleGreen, img)
	case 3:
		margin := 8*2*r - r
		x2 := singleCaptchaWidth/2 + xoffset
		x1 := x2 - margin
		x3 := x2 + margin
		y2 := captchaHeight / 2
		y1 := y2 - margin
		y3 := y2 + margin

		thick := 2 * r
		drawRing(x1, y1, 8, thick, colorOuterBlue, img)
		drawRing(x1, y1, 4, thick, colorOuterBlue, img)
		drawCircle(x1, y1, r+1, colorOuterBlue, img)

		drawRing(x2, y2, 8, thick, colorCenterRed, img)
		drawRing(x2, y2, 4, thick, colorCenterRed, img)
		drawCircle(x2, y2, r+1, colorCenterRed, img)

		drawRing(x3, y3, 8, thick, colorMiddleGreen, img)
		drawRing(x3, y3, 4, thick, colorMiddleGreen, img)
		drawCircle(x3, y3, r+1, colorMiddleGreen, img)
	case 4:
		marginx := 5 * 2 * r
		marginy := 2 * (marginx - 2)
		x := singleCaptchaWidth/2 + xoffset
		y := captchaHeight / 2
		drawNormalRingAndCircle(x-marginx, y-marginy, r, colorOuterBlue, img)

		drawNormalRingAndCircle(x-marginx, y+marginy, r, colorMiddleGreen, img)

		drawNormalRingAndCircle(x+marginx, y-marginy, r, colorMiddleGreen, img)

		drawNormalRingAndCircle(x+marginx, y+marginy, r, colorOuterBlue, img)
	case 5:
		marginx := 6 * 2 * r
		marginy := 2 * (marginx - 2*r)
		x := singleCaptchaWidth/2 + xoffset
		y := captchaHeight / 2
		drawNormalRingAndCircle(x, y, r, colorCenterRed, img)

		drawNormalRingAndCircle(x-marginx, y-marginy, r, colorOuterBlue, img)

		drawNormalRingAndCircle(x-marginx, y+marginy, r, colorMiddleGreen, img)

		drawNormalRingAndCircle(x+marginx, y-marginy, r, colorMiddleGreen, img)

		drawNormalRingAndCircle(x+marginx, y+marginy, r, colorOuterBlue, img)
	case 6:
		margin := 4*2*r + r
		x := singleCaptchaWidth/2 + xoffset
		y := captchaHeight / 2
		drawNormalRingAndCircle(x-margin, y-2*margin, r, colorMiddleGreen, img)

		drawNormalRingAndCircle(x+margin, y-2*margin, r, colorMiddleGreen, img)

		// next four
		// red

		drawNormalRingAndCircle(x-margin, y+margin, r, colorCenterRed, img)

		drawNormalRingAndCircle(x-margin, y+3*margin, r, colorCenterRed, img)

		drawNormalRingAndCircle(x+margin, y+margin, r, colorCenterRed, img)

		drawNormalRingAndCircle(x+margin, y+3*margin, r, colorCenterRed, img)
	case 7:
		margin := 4 * 2 * r
		x := singleCaptchaWidth/2 + xoffset
		y := captchaHeight / 2

		drawNormalRingAndCircle(x-2*margin, y-3*margin+3*r, r, colorMiddleGreen, img)

		drawNormalRingAndCircle(x, y-2*margin+r, r, colorMiddleGreen, img)

		drawNormalRingAndCircle(x+2*margin, y-margin-3*r, r, colorMiddleGreen, img)

		// next four
		// red

		drawNormalRingAndCircle(x-margin, y+margin, r, colorCenterRed, img)

		drawNormalRingAndCircle(x-margin, y+3*margin, r, colorCenterRed, img)

		drawNormalRingAndCircle(x+margin, y+margin, r, colorCenterRed, img)

		drawNormalRingAndCircle(x+margin, y+3*margin, r, colorCenterRed, img)

	case 8:
		margin := 4*2*r + r
		x := singleCaptchaWidth/2 + xoffset
		y := captchaHeight / 2

		drawNormalRingAndCircle(x-margin, y-3*margin, r, colorOuterBlue, img)
		drawNormalRingAndCircle(x-margin, y+margin, r, colorOuterBlue, img)
		drawNormalRingAndCircle(x-margin, y+3*margin, r, colorOuterBlue, img)

		drawNormalRingAndCircle(x-margin, y-margin, r, colorOuterBlue, img)
		drawNormalRingAndCircle(x+margin, y-margin, r, colorOuterBlue, img)

		drawNormalRingAndCircle(x+margin, y-3*margin, r, colorOuterBlue, img)
		drawNormalRingAndCircle(x+margin, y+margin, r, colorOuterBlue, img)
		drawNormalRingAndCircle(x+margin, y+3*margin, r, colorOuterBlue, img)
	case 9:
		marginx := 7 * 2 * r
		marginy := 2*marginx - r
		x := singleCaptchaWidth/2 + xoffset
		y := captchaHeight / 2

		drawNormalRingAndCircle(x, y-marginy, r, colorOuterBlue, img)
		drawNormalRingAndCircle(x-marginx, y-marginy, r, colorOuterBlue, img)
		drawNormalRingAndCircle(x+marginx, y-marginy, r, colorOuterBlue, img)

		drawNormalRingAndCircle(x, y, r, colorCenterRed, img)
		drawNormalRingAndCircle(x-marginx, y, r, colorCenterRed, img)
		drawNormalRingAndCircle(x+marginx, y, r, colorCenterRed, img)

		drawNormalRingAndCircle(x, y+marginy, r, colorMiddleGreen, img)
		drawNormalRingAndCircle(x-marginx, y+marginy, r, colorMiddleGreen, img)
		drawNormalRingAndCircle(x+marginx, y+marginy, r, colorMiddleGreen, img)
	}
}
