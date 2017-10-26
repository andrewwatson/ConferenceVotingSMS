package voteconf

import (
	"net/http"
	"time"
)

func dashHandler(rw http.ResponseWriter, r *http.Request) {
	if requireAuthentication(rw, r) == nil {
		return
	}

	data := struct {
		Title    string
		Subtitle string
		Year     int
	}{
		"Conference Session Rater",
		"Welcome Screen",
		time.Now().Year(),
	}

	err := templates.ExecuteTemplate(rw, "dashboard", &data)
	if err != nil {
		handleError(rw, r, err)
		return
	}
}
