package gtea

import (
	"fmt"
)

// Background start bg job
func (app *application) Background(fn func()) {
	rcv := func() {
		if r := recover(); r != nil {
			app.Logger.PrintErr(fmt.Errorf("%s", r), nil)
		}
	}
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		defer rcv()
		fn()
	}()
}
