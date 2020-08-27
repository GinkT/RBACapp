package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/mikespook/gorbac"

	abclientstate "github.com/volatiletech/authboss-clientstate"
	abrenderer "github.com/volatiletech/authboss-renderer"
	"github.com/volatiletech/authboss/v3"
	_ "github.com/volatiletech/authboss/v3/auth"
	"github.com/volatiletech/authboss/v3/defaults"

)

const (
	sessionCookieName = "ginkt"
)

var (
	RBAC            = gorbac.New()
	RBACPermissions = make(gorbac.Permissions)

	ab = authboss.New()
	database = NewMemStorer()
	sessionStore abclientstate.SessionStorer
	cookieStore  abclientstate.CookieStorer
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, there!")
}

func RBACInit() {
	rUser := gorbac.NewStdRole("user")
	rAdmin := gorbac.NewStdRole("admin")

	RBACPermissions["/foo"] = gorbac.NewStdPermission("/foo")
	RBACPermissions["/bar"] = gorbac.NewStdPermission("/bar")
	RBACPermissions["/sigma"] = gorbac.NewStdPermission("/sigma")

	rUser.Assign(RBACPermissions["/foo"])
	rUser.Assign(RBACPermissions["/bar"])
	rAdmin.Assign(RBACPermissions["/foo"])
	rAdmin.Assign(RBACPermissions["/bar"])
	rAdmin.Assign(RBACPermissions["/sigma"])

	RBAC.Add(rUser)
	RBAC.Add(rAdmin)
}

func RBACMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing RBAC middleware")

		user, err := ab.CurrentUser(r)
		if err != nil {
			log.Println("Error trying to get Current User:", err)
		}

		role := database.Users[user.GetPID()].Role
		if !RBAC.IsGranted(role, RBACPermissions[r.URL.Path], nil) {
			log.Printf("[RBAC] Not enough righs for user %s to access %s", database.Users[user.GetPID()].Name, r.URL.Path)
			fmt.Fprintf(w, "Not enough rights!")
			return
		}
		log.Printf("User %s have enought rigts to enter %s", database.Users[user.GetPID()].Name, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func authbossInit() {
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
}

func main() {
	mux := chi.NewRouter()
	mux.Use(ab.LoadClientStateMiddleware)

	RBACInit()
	authbossInit()

	if err := ab.Init(); err != nil {
		panic(err)
	}

	// Auth routes
	mux.Group(func(mux chi.Router) {
		mux.Use(authboss.Middleware2(ab, authboss.RequireNone, authboss.RespondRedirect), RBACMiddleware)
		mux.Get( "/foo", testHandler)
		mux.Get( "/bar", testHandler)
		mux.Get( "/sigma", testHandler)
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