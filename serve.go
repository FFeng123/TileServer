package main

import (
	"fmt"
	"net/http"
	"strings"
	"tileServer/db"
)

func serve(port int, dbf string) error {
	if port == 0 {
		port = 8000
	}

	dbmgr, err := db.NewReader(dbf)
	if err != nil {
		return err
	}

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := strings.Split(r.URL.Path, "/")
			if len(s) != 4 {
				w.WriteHeader(http.StatusNotFound)
			}
			marker := db.NewMarker(s[2], s[3], s[1])
			data, err := dbmgr.Read(marker)

			if err != nil {
				fmt.Printf("(%s) Error: %s\n", marker.String(), err.Error())
				w.WriteHeader(http.StatusNotFound)
				return
			}
			fmt.Printf("(%s) OK\n", marker.String())
			w.Write(data)

		}),
		Addr: fmt.Sprintf(":%d", port),
	}

	return server.ListenAndServe()
}
