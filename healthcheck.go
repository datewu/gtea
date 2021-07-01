package gtea

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"envirionment": app.config.env,
		},
	}
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}
}
