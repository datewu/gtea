package captcha

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestNewPNG123456(t *testing.T) {
	out, err := os.Create("captcha_test.png")
	if err != nil {
		log.Fatalln(err)
		return
	}
	out.Write(genPNG(123456))
}
func TestNewPNG789033(t *testing.T) {
	out, err := os.Create("captcha_test.png")
	if err != nil {
		log.Fatalln(err)
		return
	}
	out.Write(genPNG(789033))
}
func TestNewPNG393938(t *testing.T) {
	out, err := os.Create("captcha_test.png")
	if err != nil {
		log.Fatalln(err)
		return
	}
	out.Write(genPNG(393938))
}

func TestDigits(t *testing.T) {
	p := 6
	n := randomN(p)
	fmt.Println("number", n)
	ds := digits(n, p)
	fmt.Println("digits:", ds)
}
func TestNewCaptcha(t *testing.T) {
	c := NewCaptcha()
	fmt.Printf("%#v\n", c)
	out, err := os.Create("captcha_test.png")
	if err != nil {
		log.Fatalln(err)
		return
	}
	out.Write(genPNG(c.Tag))
	fmt.Println(c.PNG())
}

func TestRand(t *testing.T) {
	fmt.Println(randomN(2))
	fmt.Println(randomN(3))
	fmt.Println(randomN(4))
	fmt.Println(randomN(5))
	fmt.Println(randomN(6))
}
