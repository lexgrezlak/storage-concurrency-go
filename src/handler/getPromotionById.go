package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"storage-api/src/service"
)

const (
	ID = "id"
)

func GetPromotionById(datastore service.PromotionDatastore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		vars := mux.Vars(r)
		id := vars[ID]

		// Get the promotion (from redis)
		promotion, err := datastore.GetPromotionById(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		// Prepare the response data.
		payload, err := json.Marshal(promotion)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}
