package main

import (
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"task4/mydatabase"
	"task4/handlers"
)

func main() {
	mydatabase.Initdb()

	router := mux.NewRouter()

	router.HandleFunc("/apitrello/project", handlers.GetProjectsHandler).Methods("GET")
	router.HandleFunc("/apitrello/project/{id}", handlers.GetProjectHandler).Methods("GET")
	router.HandleFunc("/apitrello/project", handlers.PostProjectHandler).Methods("POST")
	router.HandleFunc("/apitrello/project", handlers.PutProjectHander).Methods("PUT")
	router.HandleFunc("/apitrello/project/{id}", handlers.DeleteProjectHander).Methods("DELETE")

	router.HandleFunc("/apitrello/column/{id}", handlers.GetColumnHandler).Methods("GET")
	router.HandleFunc("/apitrello/column", handlers.PostColumnHandler).Methods("POST")
	router.HandleFunc("/apitrello/column/newposition", handlers.UpdateColumnPositionHandler).Methods("POST")
	router.HandleFunc("/apitrello/column/{id}", handlers.DeleteColumnHandler).Methods("DELETE")
	
	router.HandleFunc("/apitrello/task/{id}", handlers.GetTaskHandler).Methods("GET")
	router.HandleFunc("/apitrello/task", handlers.PostTaskHandler).Methods("POST")
	router.HandleFunc("/apitrello/task/newposition", handlers.UpdateTaskPosition).Methods("POST")
	router.HandleFunc("/apitrello/task/{id}", handlers.DeleteTaskHandler).Methods("DELETE")

	router.HandleFunc("/apitrello/comment", handlers.PostCommentHandler).Methods("POST")
	router.HandleFunc("/apitrello/comment/{id}", handlers.GetCommentHandler).Methods("GET")
	router.HandleFunc("/apitrello/comment/{id}", handlers.DeleteCommentHandler).Methods("DELETE")


	http.ListenAndServe(":8000", router)
}