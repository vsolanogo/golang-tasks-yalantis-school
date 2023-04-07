package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"encoding/json"
	"task4/mydatabase"
	"task4/models"
)

func GetColumnHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) 

	res := mydatabase.GetProjectColumn(params["id"])

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res);	
}

func PostColumnHandler(w http.ResponseWriter, r *http.Request) {
	var newColumn models.ProjectColumn
	
	err := json.NewDecoder(r.Body).Decode(&newColumn)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	res := mydatabase.PostProjectColumn(newColumn)

	json.NewEncoder(w).Encode(res);	
}

func UpdateColumnPositionHandler(w http.ResponseWriter, r *http.Request) {
	var newColumn models.ProjectColumn
	
	err := json.NewDecoder(r.Body).Decode(&newColumn)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	mydatabase.UpdateColumnPosition(newColumn)

	json.NewEncoder(w);	
}

func DeleteColumnHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) 
	
	mydatabase.DeleteColumn(params["id"])

	resp := make(map[string]string)
	resp["message"] = "Status OK"

	json.NewEncoder(w).Encode(resp);	
}