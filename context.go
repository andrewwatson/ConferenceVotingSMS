package voteconf

import (
	"golang.org/x/net/context"

	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

/*
	HostContext returns a context object with a value added for the host header
	This is very useful for generating a CSP nonce and making it available to
	all your handlers, for example.

	Detecting the bit.ly agent is important if you're trying to tell if people
	opened a link you sent them via SMS/Email

	Also, this function takes care of setting namespaces in your App Engine context
	So that you can pass this context to Datastore calls
*/
func hostContext(r *http.Request) context.Context {
	ctx := appengine.NewContext(r)
	bitlyAgent := false
	if r.UserAgent() == "bitlybot/3.0 (+http://bit.ly/)" {
		bitlyAgent = true
	}

	nonce := r.Header.Get("CSP-Nonce")

	ctx = context.WithValue(ctx, contextFlagBitlyAgent, bitlyAgent)
	ctx = context.WithValue(ctx, contextFlagCSPNonce, nonce)
	ctx = context.WithValue(ctx, contextFlagFeatureFlags, featureFlags)

	var err error
	if featureFlags[featureFlagNamespace] && r.Host != defaultHostName {
		log.Infof(ctx, "Host is not default, setting namespace to %s", r.Host)
		ctx, err = appengine.Namespace(ctx, r.Host)
		if err != nil {
			log.Criticalf(ctx, "Namespace Error: %s", err.Error())
		}
	}

	return context.WithValue(ctx, contextFlagHost, r.Host)
}
