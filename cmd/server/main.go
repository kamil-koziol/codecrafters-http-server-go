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
		http.WriteResponse(r, w, http.StatusOK, nil, nil)
	})

	router.GET("/echo/{str}", func(r *http.Request, w io.Writer) {
		str := r.GetPath("str")
		http.WriteResponse(r, w, http.StatusOK, []byte(str), nil)
	})

	router.GET("/user-agent", func(r *http.Request, w io.Writer) {
		userAgent, _ := r.Headers.Get("User-Agent")

		http.WriteResponse(r, w, http.StatusOK, []byte(userAgent), nil)
	})

	router.GET("/files/*", func(r *http.Request, w io.Writer) {
		path := r.GetPath("*")

		fpath := filepath.Join(*directory, path)
		f, err := os.Open(fpath)
		if err != nil {
			if os.IsNotExist(err) {
				http.WriteResponse(r, w, http.StatusNotFound, nil, nil)
				return
			}

			http.WriteResponse(r, w, http.StatusInternalServerError, nil, nil)
			return
		}
		defer f.Close()

		fileContents, err := io.ReadAll(f)
		if err != nil {
			http.WriteResponse(r, w, http.StatusInternalServerError, nil, nil)
			return
		}

		h := http.NewHeaders()
		h.Set("Content-Type", "application/octet-stream")

		http.WriteResponse(r, w, http.StatusOK, fileContents, h)
	})

	router.POST("/files/{filename}", func(r *http.Request, w io.Writer) {
		filename := r.GetPath("filename")

		fpath := filepath.Join(*directory, filename)
		if err := os.WriteFile(fpath, r.Body, 0644); err != nil {
			http.WriteResponse(r, w, http.StatusInternalServerError, nil, nil)
			return
		}

		http.WriteResponse(r, w, http.StatusCreated, nil, nil)
	})

	srv := http.Server{
		Router: router,
	}
	err := srv.Run("0.0.0.0:4221")
	if err != nil {
		log.Fatalf("failure during run: %v", err)
	}
}
