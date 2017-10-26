package voteconf

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func requireAuthentication(rw http.ResponseWriter, r *http.Request) *user.User {
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	if u == nil {
		url, _ := user.LoginURL(ctx, r.RequestURI)
		http.Redirect(rw, r, url, http.StatusTemporaryRedirect)
		return nil
	}

	return u
}

func truncateString(input string, length int) (string, error) {

	inputLength := len(input)

	shorter := input
	if inputLength > length-1 {
		shorter = string([]byte(input)[0:length]) + "..."
	}

	return shorter, nil
}
