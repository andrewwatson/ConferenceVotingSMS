package voteconf

import (
	// "fmt"
	"github.com/gorilla/mux"
	// "html"
	"appengine"
	"appengine/datastore"
	"fmt"
	"net/http"

	"strings"
)

type VoteStorage struct {
	PhoneNumber string
	Hashtag     string
	Rating      int
	Comment     string
}

// Because App Engine owns main and starts the HTTP service,
// we do our setup during initialization.
func init() {

	r := mux.NewRouter()

	// r.HandleFunc("/", HomeHandler)
	r.Path("/sms").Methods("POST").HandlerFunc(SMSHandler)
	r.Path("/report").Methods("GET").HandlerFunc(ReportHandler)

	r.Path("")
	http.Handle("/", r)
}

func ReportHandler(rw http.ResponseWriter, req *http.Request) {
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

func SMSHandler(rw http.ResponseWriter, req *http.Request) {

	c := appengine.NewContext(req)

	message := req.PostFormValue("Body")
	rw.Header().Add("Content-type", "text/xml")

	segments := strings.SplitN(message, " ", 3)
	segmentCount := len(segments)

	twiml := ""

	if segmentCount == 3 {

		voteValue, err := validateVote(segments[1])

		if err != nil {
			twiml = fmt.Sprintf("<Response><Message>%s</Message></Response>", err.Error())
		} else {
			twiml = "<Response><Message>Thank you for your feedback. Enjoy the rest of #CSSDevConf!</Message></Response>"

			comment := fmt.Sprintf("'%s'", segments[2])

			err = validateHashtag(segments[0])

			if err != nil {
				twiml = fmt.Sprintf("<Response><Message>%s</Message></Response>", err.Error())

			} else {

				phoneNumber := req.PostFormValue("From")
				v := VoteStorage{phoneNumber, segments[0], voteValue, comment}
				_, err = datastore.Put(c, datastore.NewIncompleteKey(c, "vote", nil), &v)

				if err != nil {
					twiml = "<Response><Message>An Error Happened While Recording Your Feedback.  We're on it.</Message></Response>"
				}

			}

		}

	} else {
		twiml = `<Response><Message>Your Message was not formatted properly. Sad trombone.</Message></Response>`

	}

	rw.Write([]byte(twiml))
}
