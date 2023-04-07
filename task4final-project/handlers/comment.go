package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"encoding/json"
	"task4/mydatabase"
	"task4/models"
)

func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	var i models.Comment
	
	err := json.NewDecoder(r.Body).Decode(&i)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	res := mydatabase.InsertComment(i)

	json.NewEncoder(w).Encode(res);	
}


func GetCommentHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) 

	res := mydatabase.SelectComment(params["id"])

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res);	
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) 

	mydatabase.DeleteComment(params["id"])

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w);	
}