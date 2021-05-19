package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

const (
	backURL = "4pda.ru"
	port = "8080"
	maxCountCache=2
)
func  ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "count is %d\n", )
}

func main()  {
	http.HandleFunc("/endpoint", ServeHTTP)

	log.Fatal(http.ListenAndServe(":8080", nil))

	req, err := http.NewRequest("GET", backURL, nil)

	if err != nil {
		log.Fatal(err)
	}


	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	body, err := httputil.DumpResponse(resp, true)

	if err != nil {
		log.Fatal(err)
	}

	// write to disk, or something
	// ...

	// wrap the cached response

	r := bufio.NewReader(bytes.NewReader(body))
	// ReadResponse by default assumes the request for the response was a "GET" requested
	// If you want the method to be different, you must pass an http.Request to ReadResponse (instead of nil)
	resp, err = http.ReadResponse(r, nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v", resp)
}