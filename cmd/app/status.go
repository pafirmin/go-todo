package main

import "net/http"

func (app *application) showStatus(w http.ResponseWriter, r *http.Request) {
	rsp := responsePayload{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	app.writeJSON(w, http.StatusOK, rsp)
}
