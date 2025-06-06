package internal

import (
	"encoding/json"
	"net/http"
	"time"
)

func Routes(store *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/articles", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		listHandler(store, w, r)
	})
	return mux
}

func listHandler(store *Store, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	sortAsc := q.Get("sort") == "asc"

	var fromPtr, toPtr *time.Time
	if v := q.Get("from"); v != "" {
		t, err := time.Parse("2007-11-26", v)
		if err != nil {
			http.Error(w, "bad 'from' date (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		fromPtr = &t
	}
	if v := q.Get("to"); v != "" {
		t, err := time.Parse("2007-11-26", v)
		if err != nil {
			http.Error(w, "bad 'to' date (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		toPtr = &t
	}

	arts, err := store.List(fromPtr, toPtr, sortAsc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(arts)
}
