package handler

import (
	"github.com/gorilla/mux"
	"jayant/database/dbHelper"
	"jayant/middlewares"
	"jayant/models"
	"jayant/utils"
	"net/http"
	"strconv"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserContext(r)
	body := models.Task{}
	err := utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse request body")
		return
	}
	err = dbHelper.CreateTask(user.ID, body)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create task")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AllTask(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserContext(r)
	var err error
	var isCompleted bool
	var isFiltered bool
	par := r.URL.Query().Get("choice")
	if par != "" {
		isFiltered = true
		isCompleted, err = strconv.ParseBool(par)
		if err != nil {
			utils.RespondError(w, http.StatusBadRequest, err, "input error")
			return
		}
	}
	task, err := dbHelper.AllTask(user.ID, isFiltered, isCompleted)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to retrieve tasks")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"task": task,
	})
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserContext(r)
	var body models.UpdateTask
	err := utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse request body")
		return
	}

	err = dbHelper.UpdateTask(user.ID, body)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to update todo")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func Complete(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserContext(r)
	i := mux.Vars(r)["taskId"]
	id, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Error: no characters are allowed")
		return
	}
	err = dbHelper.Complete(id, user.ID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to update")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	user := middlewares.UserContext(r)
	i := mux.Vars(r)["TaskId"]
	id, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Error: no characters are allowed")
		return
	}
	Err := dbHelper.DeleteTask(id, user.ID)
	if Err != nil {
		utils.RespondError(w, http.StatusInternalServerError, Err, "failed to delete")
		return
	}
	w.WriteHeader(http.StatusOK)
}
