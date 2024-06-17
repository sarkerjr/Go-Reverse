package proxy

import (
	"io"
	"log"
	"net"
	"net/http"

	"sarkerjr.com/go-reverse/pkg/client"
)

type Proxy struct {
	config *Config
	client *http.Client
}

func NewProxy() *Proxy {
	return &Proxy{
		config: GetConfig(),
		client: client.NewHTTPClient(),
	}
}

func (p *Proxy) HandleRequest(w http.ResponseWriter, req *http.Request) {
	// Forward the IP address of the client to the origin server
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Printf("Error parsing client IP address: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("X-Forwarded-For", host)

	// Create a new request to the origin server, copying headers from the original request
	originReq, err := http.NewRequest(req.Method, p.config.OriginURL+req.URL.Path, req.Body)
	if err != nil {
		log.Printf("Error creating request to backend server: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Copy relevant headers to the origin request
	for key, values := range req.Header {
		for _, value := range values {
			originReq.Header.Add(key, value)
		}
	}

	// Make the request to the backend server
	backendResp, err := p.client.Do(originReq)
	if err != nil {
		log.Printf("Error forwarding request to backend: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer backendResp.Body.Close()

	// Copy response headers from the backend to the client
	for key, values := range backendResp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy response body and headers from the backend to the client
	w.WriteHeader(backendResp.StatusCode)

	// Copy the response body from the backend to the client
	if _, err := io.Copy(w, backendResp.Body); err != nil {
		log.Printf("Error copying response body to client: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
