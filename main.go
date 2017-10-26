package voteconf

import (
	"fmt"
	"html"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"appengine"
	"appengine/datastore"

	"strings"
)

var (
	templates *template.Template
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

	funcMap := template.FuncMap{

		"Title":       strings.Title,
		"Shorten":     truncateString,
		"Unescape":    html.UnescapeString,
		"StringInt":   strconv.Itoa,
		"StringInt64": strconv.FormatInt,
	}
	// Preload all the templates at startup time instead of in each handler function
	// All templates must start with a {{define "name"}} block and end with {{end}}
	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

	r := mux.NewRouter()

	// Serve up static assets directly
	r.PathPrefix("/assets").Handler(http.FileServer(http.Dir(".")))

	r.HandleFunc("/", HomeHandler)

	r.Path("/h").Methods("GET").HandlerFunc(HomeHandler)
	r.Path("/sms").Methods("POST").HandlerFunc(SMSHandler)
	r.Path("/report").Methods("GET").HandlerFunc(ReportHandler)

	// r.NotFoundHandler = http.NotFoundHandler()

	http.Handle("/", r)
}

// ReportHandler generates the CSV file output
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

// SMSHandler processes incoming messages
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

// HomeHandler is the home page
func HomeHandler(rw http.ResponseWriter, r *http.Request) {

	data := struct {
		Title    string
		Subtitle string
		Year     int
	}{
		"Conference Session Rater",
		"Welcome Screen",
		time.Now().Year(),
	}

	err := templates.ExecuteTemplate(rw, "main", &data)
	if err != nil {
		handleError(rw, r, err)
		return
	}
}

func truncateString(input string, length int) (string, error) {

	inputLength := len(input)

	shorter := input
	if inputLength > length-1 {
		shorter = string([]byte(input)[0:length]) + "..."
	}

	return shorter, nil
}
