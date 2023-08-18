package captcha

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var pp Pnger = MajongPng{}

func TestNewPNG123456(t *testing.T) {
	out, err := os.Create("captcha_test.png")
	if err != nil {
		log.Fatalln(err)
		return
	}
	out.Write(pp.PNG(123456))
}
func TestNewPNG789033(t *testing.T) {
	out, err := os.Create("captcha_test.png")
	if err != nil {
		log.Fatalln(err)
		return
	}
	out.Write(pp.PNG(789033))
}
func TestNewPNG393938(t *testing.T) {
	out, err := os.Create("captcha_test.png")
	if err != nil {
		log.Fatalln(err)
		return
	}
	out.Write(pp.PNG(393938))
}

func TestDigits(t *testing.T) {
	p := 6
	n := randomN(p)
	fmt.Println("number", n)
	ds := digits(n, p)
	fmt.Println("digits:", ds)
}

func TestNewCaptcha(t *testing.T) {
	c := NewCaptcha(nil)
	fmt.Printf("%#v\n", c)
	out, err := os.Create("captcha_test.png")
	if err != nil {
		log.Fatalln(err)
		return
	}
	out.Write(pp.PNG(c.Tag))
	fmt.Println(c.PNG())
}

func TestRand(t *testing.T) {
	for i := 0; i < 7; i++ {
		fmt.Println(randomN(i))
	}
}
