package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"encoding/json"
	"task4/mydatabase"
	"task4/models"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) 

	res := mydatabase.SelectTask(params["id"])

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res);	
}

func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var i models.Task
	
	err := json.NewDecoder(r.Body).Decode(&i)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	res := mydatabase.InsertTask(i)

	json.NewEncoder(w).Encode(res);	
}

func UpdateTaskPosition(w http.ResponseWriter, r *http.Request) {
	var newTask models.Task
	
	err := json.NewDecoder(r.Body).Decode(&newTask)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	if(newTask.ProjectColumnId != 0) {
		mydatabase.UpdateTaskPosition(newTask)
	} else {
		mydatabase.UpdateTaskPositionInScopeOfTheColumn(newTask)
	}

	json.NewEncoder(w);	
}
	
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) 

	mydatabase.DeleteTask(params["id"])

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w);	
}