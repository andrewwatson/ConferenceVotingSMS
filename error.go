package voteconf

import (
	"net/http"
)

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	// 	w.WriteHeader(http.StatusInternalServerError)

	data := struct {
		Error    string
		Title    string
		LoggedIn bool
	}{
		err.Error(),
		"Error",
		false,
	}

	templates.ExecuteTemplate(w, "error", data)

	return
}
