package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	proxy    *httputil.ReverseProxy
	upstream *url.URL
}


func NewReverseProxy(upstreamURL string) (*ReverseProxy, error) {
	upstream, err := url.Parse(upstreamURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(upstream)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = upstream.Host

		req.Header.Set("X-Sidecar-Version", "1.0")

		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Real-IP", req.RemoteAddr)

		log.Printf("=== Request Headers for %s %s ===", req.Method, req.URL.Path)
		for name, values := range req.Header {
			for _, value := range values {
				log.Printf("  %s: %s", name, value)
			}
		}
		log.Printf("=== End Headers ===")
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("X-Proxy", "opentelemetry-sidecar")

		resp.Header.Del("Server")

		log.Printf("Response from upstream: %d %s", resp.StatusCode, resp.Status)
		return nil
	}

	return &ReverseProxy{
		proxy:    proxy,
		upstream: upstream,
	}, nil
}


func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var countingReader *CountingReader
	if r.Body != nil {
		countingReader = NewCountingReader(r.Body)
		r.Body = countingReader
	}

	rp.proxy.ServeHTTP(w, r)

	if countingReader != nil {
		log.Printf("Request body bytes read: %d", countingReader.BytesRead())
	}
}
func (rp *ReverseProxy) Upstream() *url.URL {
	return rp.upstream
}
