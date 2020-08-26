package main

import (
	"encoding/base64"
	"github.com/gorilla/sessions"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	"github.com/volatiletech/authboss-renderer"
	"github.com/volatiletech/authboss/v3"
	_ "github.com/volatiletech/authboss/v3/auth"
	"github.com/volatiletech/authboss/v3/defaults"
	"time"

	"github.com/go-chi/chi"

	"fmt"
	"log"
	"net/http"
)

const (
	sessionCookieName = "ginkt"
)

var (
	database = NewMemStorer()
	sessionStore abclientstate.SessionStorer
	cookieStore  abclientstate.CookieStorer
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, there!")
}

func main() {
	ab := authboss.New()

	mux := chi.NewRouter()
	mux.Use(ab.LoadClientStateMiddleware)

	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
	cookieStore = abclientstate.NewCookieStorer(cookieStoreKey, nil)
	cookieStore.HTTPOnly = false
	cookieStore.Secure = false
	sessionStore = abclientstate.NewSessionStorer(sessionCookieName, sessionStoreKey, nil)
	cstore := sessionStore.Store.(*sessions.CookieStore)
	cstore.Options.HttpOnly = false
	cstore.Options.Secure = false
	cstore.MaxAge(int((30 * 24 * time.Hour) / time.Second))

	ab.Config.Storage.Server = database
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore

	ab.Config.Paths.Mount = ""
	ab.Config.Paths.RootURL = "http://localhost:8000"

	ab.Config.Core.ViewRenderer = abrenderer.NewHTML("", "")

	defaults.SetCore(&ab.Config, false, false)

	if err := ab.Init(); err != nil {
		panic(err)
	}

	// Auth routes
	mux.Group(func(mux chi.Router) {
		mux.Use(authboss.Middleware2(ab, authboss.RequireNone, authboss.RespondRedirect))
		mux.Get( "/foo", testHandler)
		mux.Get( "/bar", testHandler)
		mux.MethodFunc("GET", "/get", testHandler)
	})

	// Unauth routes
	mux.Group(func(mux chi.Router) {
		mux.Use(authboss.ModuleListMiddleware(ab))
		mux.Mount("/", http.StripPrefix("", ab.Config.Core.Router))
		mux.Get("/", testHandler)
	})

	log.Println("Starting to listen and serve on port 8000")
	log.Fatalln(http.ListenAndServe(":8000", mux))
}