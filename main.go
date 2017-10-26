package voteconf

import (
	"html"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/appengine/user"

	"google.golang.org/appengine"

	"strings"
)

var (
	templates    *template.Template
	featureFlags map[string]bool
)

type VoteStorage struct {
	PhoneNumber string
	Hashtag     string
	Rating      int
	Comment     string
}
type contextKey int

const (
	contextFlagFeatureFlags contextKey = 0
	contextFlagBitlyAgent   contextKey = 1
	contextFlagHost         contextKey = 2
	contextFlagCSPNonce     contextKey = 3
	featureFlagCSRF                    = "FEATUREFLAG_CSRF"
	featureFlagNamespace               = "FEATUREFLAG_NAMESPACE"
	defaultHostName                    = ""
	defaultTimeZone                    = "America/New_York"
)

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

	r.HandleFunc("/", homeHandler)

	r.Path("/sms").Methods("POST").HandlerFunc(messageHandler)
	r.Path("/report").Methods("GET").HandlerFunc(reportHandler)
	r.Path("/dashboard").Methods("GET").HandlerFunc(dashHandler)

	r.Path("/logout").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		logoutURL, _ := user.LogoutURL(ctx, "/")
		http.Redirect(w, r, logoutURL, http.StatusTemporaryRedirect)
		return
	})

	r.NotFoundHandler = http.NotFoundHandler()

	http.Handle("/", r)
}

// HomeHandler is the home page
func homeHandler(rw http.ResponseWriter, r *http.Request) {

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
