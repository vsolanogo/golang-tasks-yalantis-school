package handlers

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"encoding/json"
	"reflect"
	"task4/models"
	"task4/mydatabase"
)

func GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) 

	queriedProjects := mydatabase.GetProject(params["id"])

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queriedProjects);	
}

func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) 
	log.Println(reflect.TypeOf(params["id"]))

	res := mydatabase.GetProjects()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res);	
}

func PostProjectHandler(w http.ResponseWriter, r *http.Request) {
	var newProject models.Project
	
	err := json.NewDecoder(r.Body).Decode(&newProject)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	res := mydatabase.PostProject(newProject)

	json.NewEncoder(w).Encode(res);	
}

func PutProjectHander(w http.ResponseWriter, r *http.Request) {
	
	var updatedProject models.Project
	
	err := json.NewDecoder(r.Body).Decode(&updatedProject)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	err = mydatabase.PutProject(updatedProject)

	if(err != nil) {
		log.Fatal(err)
	}

	resp := make(map[string]string)
	resp["message"] = "Status OK"

	json.NewEncoder(w).Encode(resp);	
}

func DeleteProjectHander(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) 
	err := mydatabase.DeleteProject(params["id"])
	
	if(err != nil) {
		log.Fatal(err)
	}

	resp := make(map[string]string)
	resp["message"] = "Status OK"

	json.NewEncoder(w).Encode(resp);	
}