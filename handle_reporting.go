package voteconf

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// ReportHandler generates the CSV file output
func reportHandler(rw http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)

	var Votes []VoteStorage
	q := datastore.NewQuery("vote")
	_, err := q.GetAll(c, &Votes)

	if err != nil {

	}

	rw.Header().Add("Content-type", "text/plain")

	for _, vote := range Votes {
		rw.Write([]byte(fmt.Sprintf("%s, %s,%d,%s\n", vote.PhoneNumber, vote.Hashtag, vote.Rating, vote.Comment)))
	}

}
