package gutil

import "net/http"

// GetIP gets the provided requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func IP(r *http.Request) string {
	if r == nil {
		return ""
	}

	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
