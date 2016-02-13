// Package handlers provides HTTP request handlers.
package handlers

import (
	"net/http"

	"github.com/thomaso-mirodin/shawty/storages"
)

func EncodeHandler(storage storages.Storage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		if url := r.PostFormValue("url"); url != "" {
			code, err := storage.Save(url)
			if err != nil {
				http.Error(w, err, 500)
				return
			}
			w.Write([]byte(code))
		}
	}

	return http.HandlerFunc(handleFunc)
}

func DecodeHandler(storage storages.Storage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Path[len("/dec/"):]

		url, err := storage.Load(code)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("URL Not Found. Error: " + err.Error() + "\n"))
			return
		}

		w.Write([]byte(url))
	}

	return http.HandlerFunc(handleFunc)
}

func RedirectHandler(storage storages.Storage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Path[len("/red/"):]

		url, err := storage.Load(code)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("URL Not Found. Error: " + err.Error() + "\n"))
			return
		}

		http.Redirect(w, r, string(url), 301)
	}

	return http.HandlerFunc(handleFunc)
}
