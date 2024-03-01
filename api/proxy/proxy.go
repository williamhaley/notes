package proxy

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Options struct {
	DestinationOrigin string
	StripPrefix       string
	AdditionalHeaders map[string]string
}

func ProxyTo(opts *Options) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		url := fmt.Sprintf("%s%s", opts.DestinationOrigin, strings.TrimPrefix(req.URL.String(), opts.StripPrefix))

		req, err := http.NewRequest(req.Method, url, req.Body)

		// TODO
		req.Header.Add("Content-Type", "application/json")

		for k, v := range opts.AdditionalHeaders {
			req.Header.Add(k, v)
		}

		if err != nil {
			panic(err)
		}

		resp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()
		copyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
