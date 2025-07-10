package router

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/rsys-speerzad/stackgen/pkg/api"
	"github.com/rsys-speerzad/stackgen/pkg/api/events"
	"github.com/rsys-speerzad/stackgen/pkg/api/users"
	"github.com/rsys-speerzad/stackgen/pkg/store"
)

func NewServer() *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	return &http.Server{Addr: "localhost:" + port, Handler: newHandler()}
}

func newHandler() http.Handler {
	// initialize the router
	r := httprouter.New()
	// add user routes
	users.InitializeRouter(r, store.GetDB())
	// add event routes
	events.InitializeRouter(r, store.GetDB())
	// add gloabal options
	r.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			header := w.Header()
			header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		}
		w.WriteHeader(http.StatusNoContent)
	})
	// add global error handlers
	r.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		log.Printf("panic: %+v", err)
		api.Error(w, r, fmt.Errorf("whoops! My handler has run into a panic"), http.StatusInternalServerError)
	}
	// add method not allowed and not found handlers
	r.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.Error(w, r, fmt.Errorf("we have OPTIONS for youm but %v is not among them", r.Method), http.StatusMethodNotAllowed)
	})
	// add not found handler
	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.Error(w, r, fmt.Errorf("whatever route you've been looking for, it's not here"), http.StatusNotFound)
	})
	return r
}
