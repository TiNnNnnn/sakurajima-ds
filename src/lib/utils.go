package utils

import "net/http"

//从client的http请求头中取出object hash
func GetHashFromHeader(h http.Header) string {
	digest := h.Get("didest")
	if len(digest) < 9 {
		return ""
	}

	if digest[:8] != "SHA-256=" {
		return ""
	}

	return digest[8:]
}


