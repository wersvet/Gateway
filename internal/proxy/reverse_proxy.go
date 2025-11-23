package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewReverseProxy configures a reverse proxy targeting the provided backend URL.
func NewReverseProxy(target string) (*httputil.ReverseProxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host

		// Прокидываем user info
		if userID := req.Header.Get("X-User-ID"); userID != "" {
			req.Header.Set("X-User-ID", userID)
		}
		if username := req.Header.Get("X-Username"); username != "" {
			req.Header.Set("X-Username", username)
		}
	}

	return proxy, nil
}
