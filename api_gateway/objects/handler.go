package objects

import "net/http"

var PeersMap = []string{
	"127.0.0.1:10020",
	"127.0.0.1:10021",
	"127.0.0.1:10022",
}

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		put(w, r)
		return
	}
	if m == http.MethodGet {
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
