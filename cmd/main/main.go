package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"sarkerjr.com/go-reverse/pkg/initializer"
)

func main() {
	initializer.Initialize()

	http.HandleFunc("/", handleRequest)

	port := os.Getenv("PORT")

	log.Printf("Starting reverse proxy server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	originUrl := os.Getenv("ORIGIN_SERVER_URL")
	if originUrl == "" {
		log.Fatalln("ORIGIN_SERVER_URL environment variable is not set")
	}

	// forward the IP address of the client to the origin server
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Fatalln("Error parsing client IP address: ", err)
	}
	req.Header.Set("X-Forwarded-For", host)

	// create a new request to the origin server, copying headers from the original request
	originReq, err := http.NewRequest(req.Method, originUrl+req.URL.Path, req.Body)
	if err != nil {
		log.Fatalln("Error creating request to backend server: ", err)
	}

	// copy relevant headers to the origin request
	for key, value := range req.Header {
		originReq.Header[key] = value
	}

	// make the request to the backend server
	backendResp, err := http.DefaultClient.Do(originReq)
	if err != nil {
		fmt.Fprintf(w, "Error forwarding request to backend: %v", err)
		return
	}
	defer backendResp.Body.Close()

	// copy response headers from the backend to the client
	for key, value := range backendResp.Header {
		w.Header()[key] = value
	}

	// copy response body and headers from the backend to the client
	w.WriteHeader(backendResp.StatusCode)

	// Copy the response body from the backend to the client
	io.Copy(w, backendResp.Body)
}
