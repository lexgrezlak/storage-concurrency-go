package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"storage-api/src/service"
)

func GetPromotionById(datastore service.PromotionDatastore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		promotion, err := datastore.GetPromotionById(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		payload, err := json.Marshal(promotion)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}
}