package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

const (
	defaultConfigFile = "config.json"

	zoomAuthorizeURL = "https://devep.zoomdev.us/oauth/authorize"
	zoomTokenURL     = "https://devep.zoomdev.us/oauth/token"
)

var scopes = []string{"user:read", "meeting:read", "recording:read"}
var cfg *Config
var oauthConfig *oauth2.Config

func echoServer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	w.Header().Set("Content-Type", "text/plain")

	for k, v := range req.Header {
		fmt.Fprintf(w, "%s: %s\n", k, v)
	}

	zoomContextStr := req.Header.Get("X-Zoom-App-Context")

	if len(zoomContextStr) != 0 {
		fmt.Fprintf(w, "\n\n")
		fmt.Fprintf(w, "ctx: %s\n", zoomContextStr)
		zoomMeetingContext, err := decryptZoomContext(cfg.ClientSecret, zoomContextStr+"==")

		if err != nil {
			fmt.Fprintf(w, "error decoding meeting context: %s\n", err)
		} else {
			fmt.Fprintf(w, "decoded meeting context: %+v\n", zoomMeetingContext)
		}
	}

	fmt.Fprintf(w, "\n\nGreetings\n")

	for k, v := range vars {
		fmt.Fprintf(w, "%s: %s", k, v)
	}
}

func xxxmain() {
	r := mux.NewRouter()

	r.HandleFunc("/redirect", zoomAuthRedirectHandler).Queries("code", "{code}")
	r.HandleFunc("/", echoServer)

	http.Handle("/", r)

	err := http.ListenAndServeTLS(":9090", "ssl.dev-coursera.crt", "ssl.dev-coursera.key", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	var err error

	cfg, err = loadConfig(defaultConfigFile)

	if err != nil {
		log.Fatal(err)
	}

	oauthConfig = &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  zoomAuthorizeURL,
			TokenURL: zoomTokenURL,
		},
		RedirectURL: "",
		Scopes:      scopes,
	}

	r := mux.NewRouter()

	r.HandleFunc("/redirect", zoomAuthRedirectHandler).Queries("code", "{code}")
	r.HandleFunc("/", echoServer)

	http.Handle("/", handlers.LoggingHandler(os.Stdout, r))

	err = http.ListenAndServeTLS(":9090", "ssl.dev-coursera.crt", "ssl.dev-coursera.key", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
