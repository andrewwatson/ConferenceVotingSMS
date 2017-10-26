package voteconf

import (
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// SMSHandler processes incoming messages
func messageHandler(rw http.ResponseWriter, req *http.Request) {

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
