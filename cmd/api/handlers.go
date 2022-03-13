package main

import "net/http"

func (app *application) getTasks(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get all tasks"))
}
func (app *application) getOneTask(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get one task"))
}
func (app *application) createTask(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a task"))
}
func (app *application) updateTask(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update a task"))
}
func (app *application) deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete a task"))
}
