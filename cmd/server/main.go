package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/http-server-starter-go/internal/http"
)

func main() {
	directory := flag.String("directory", "", "Directory where the static files are stored")
	flag.Parse()

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

	router.GET("/files/*", func(r *http.Request, w io.Writer) {
		path := r.GetPath("*")

		fpath := filepath.Join(*directory, path)
		f, err := os.Open(fpath)
		if err != nil {
			if os.IsNotExist(err) {
				http.WriteResponse(w, http.StatusNotFound, nil, http.Headers{})
				return
			}

			http.WriteResponse(w, http.StatusInternalServerError, nil, http.Headers{})
			return
		}
		defer f.Close()

		fileContents, err := io.ReadAll(f)
		if err != nil {
			http.WriteResponse(w, http.StatusInternalServerError, nil, http.Headers{})
			return
		}

		h := http.Headers{}
		h.Set("Content-Type", "application/octet-stream")

		http.WriteResponse(w, http.StatusOK, fileContents, h)
	})

	srv := http.Server{
		Router: router,
	}
	err := srv.Run("0.0.0.0:4221")
	if err != nil {
		log.Fatalf("failure during run: %v", err)
	}
}
