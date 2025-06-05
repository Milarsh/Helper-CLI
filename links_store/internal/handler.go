package internal

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func Routes(store *Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/links", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listHandler(store, w)
		case http.MethodPost:
			createHandler(store, w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/v1/links/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.URL.Path[len("/api/v1/links/"):], 10, 64)
		if err != nil || id <= 0 {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			getHandler(store, w, r, id)
		case http.MethodDelete:
			deleteHandler(store, w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return mux
}

func listHandler(store *Store, w http.ResponseWriter) {
	data, err := store.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}

type createReq struct {
	URL   string `json:"url"`
	Label string `json:"label"`
}

func createHandler(store *Store, w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}
	if req.URL == "" || req.Label == "" {
		http.Error(w, "url and label required", http.StatusBadRequest)
		return
	}

	id, err := store.Add(Link{URL: req.URL, Label: req.Label})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ID int64 `json:"id"`
	}{ID: id})
}

func getHandler(store *Store, w http.ResponseWriter, r *http.Request, id int64) {
	l, found, err := store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !found {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(l)
}

func deleteHandler(store *Store, w http.ResponseWriter, r *http.Request, id int64) {
	if err := store.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
