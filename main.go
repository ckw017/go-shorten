package main

import (
	"log"
	"net"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/jessevdk/go-flags"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ckw017/go-shorten/handlers"
	"github.com/ckw017/go-shorten/storage"
)

var opts Options

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		return
	}

	store, err := createStorageFromOption(&opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Storage successfully created")

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("static")),
	)

	r := httprouter.New()
	r.Handler("GET", "/healthcheck", handlers.Healthcheck(store, "/healthcheck"))

	// Serve the index
	indexPage, err := handlers.NewIndex("static/templates/index.tmpl")
	if err != nil {
		log.Fatal("Failed to create index Page", err)
	}

	// If we don't have any matches, serve the respective go link
	r.HandleMethodNotAllowed = false
	r.NotFound = handlers.GetShort(store, indexPage)

	// Go Endpoints
	r.Handler("GET", "/go", handlers.ServeGoDashboard())

	// API handlers
	r.Handler("POST", "/", handlers.SetShort(store)) // TODO(@thomas): move this to a stable API endpoint
	if ss, ok := store.(storage.SearchableStorage); ok {
		r.Handler("GET", "/_api/v1/search", handlers.Search(ss))
	}
	if tns, ok := store.(storage.TopN); ok {
		r.Handler("GET", "/_api/v1/top_n", handlers.TopN(tns))
	}

	n.UseHandler(r)

	go func() {
		log.Printf("Starting prometheus HTTP Listener on %s", net.JoinHostPort(opts.BindHost, "8081"))
		err := http.ListenAndServe(net.JoinHostPort(opts.BindHost, "8081"), promhttp.Handler())
		if err != nil {
			log.Println(err)
		}
	}()

	log.Printf("Starting HTTP Listener on: %s", net.JoinHostPort(opts.BindHost, opts.BindPort))
	err = http.ListenAndServe(net.JoinHostPort(opts.BindHost, opts.BindPort), n)
	if err != nil {
		log.Fatal(err)
	}
}
