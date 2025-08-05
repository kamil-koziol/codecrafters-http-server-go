package main

import (
	"io"
	"log"

	"github.com/codecrafters-io/http-server-starter-go/internal/http"
)

func main() {
	router := http.Router{}

	router.GET("/", func(r *http.Request, w io.Writer) {
		http.WriteResponse(w, http.StatusOK, nil, http.Headers{})
	})

	router.GET("/echo/{str}", func(r *http.Request, w io.Writer) {
		str := r.GetPath("str")
		h := http.Headers{}
		h.Set("Content-Type", "text/plain")
		http.WriteResponse(w, http.StatusOK, []byte(str), h)
	})

	router.GET("/user-agent", func(r *http.Request, w io.Writer) {
		userAgent, _ := r.Headers.Get("User-Agent")

		h := http.Headers{}
		h.Set("Content-Type", "text/plain")

		http.WriteResponse(w, http.StatusOK, []byte(userAgent), h)
	})

	srv := http.Server{
		Router: router,
	}
	err := srv.Run("0.0.0.0:4221")
	if err != nil {
		log.Fatalf("failure during run: %v", err)
	}
}
