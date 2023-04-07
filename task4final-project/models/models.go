package models

import (
	"time"
)

type Comment struct {
	CommentId int `json:"commentid"`
	TaskId int `json:"taskid"`
	Text string `json:"text,omitempty"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

type Task struct {
	TaskId int `json:"taskid"`
	ProjectColumnId int `json:"projectcolumnid,omitempty"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	Position int `json:"position"`
	Comments []Comment `json:"comments,omitempty"`
}

type ProjectColumn struct {
	ProjectColumnId int `json:"projectcolumnid,omitempty"`
	ProjectId int `json:"projectid,omitempty"`
	Name string `json:"name,omitempty"`
	Position int `json:"position,omitempty"`
	Tasks []Task `json:"tasks"`
}

type Project struct {
	ID int `json:"projectid,omitempty"`
	Name string `json:"name"`
	Description string `json:"description"`
	ProjectColumns []ProjectColumn `json:"projectColumns,omitempty"`
}

