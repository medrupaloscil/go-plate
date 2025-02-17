package controllers

import (
	"boilerplate/models"
	"boilerplate/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(400, "Failed to read request body", w)
		return
	}
	
	user := new(models.User)
	user.LastOnline = time.Now()

	if err := json.Unmarshal(body, &user); err != nil {
		SendError(422, "fail_to_parse_body", w)
		return
	}

	password := new(models.Password)
	if err := json.Unmarshal(body, &password); err != nil {
		fmt.Println(err)
		SendError(422, "Failed to retrieve password from body", w)
		return
	}

	errs := models.CreateUser(user, password)
	if errs != nil {
		SendError(400, errs[0], w)
		return
	}

	SendResponse(user, w)
}

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(400, "Failed to read request body", w)
		return
	}
	
	user := new(models.User)
	if err := json.Unmarshal(body, &user); err != nil {
		SendError(422, "fail_to_parse_body", w)
		return
	}

	password := new(models.Password)
	if err := json.Unmarshal(body, &password); err != nil {
		fmt.Println(err)
		SendError(422, "Failed to retrieve password from body", w)
		return
	}

	if ok, logErr := models.AreLogInfosCorrect(user, password.Password); !ok {
		SendError(400, logErr, w)
		return
	}

	token, err := services.GenerateToken(user.ID)

	if err != nil {
		SendError(400, err.Error(), w)
		return
	}

	SendResponse(token, w)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 32)
	if err != nil {
		page = 0
	}

	itemsPerPage, err := strconv.ParseUint(os.Getenv("ITEMS_PER_PAGE"), 10, 32)
	if err != nil {
		itemsPerPage = 10
	}

	users, err := models.GetAllUser(int(page), int(itemsPerPage))
	if err != nil {
		SendError(400, "Failed to retrieve users", w)
		return
	}

	SendResponse(users, w)
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetUser(r.Context().Value("UserId").(uint))
	if err != nil {
		SendError(400, err.Error(), w)
	}

	SendResponse(user, w)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusBadRequest)
		return
	}

	user, err := models.GetUser(uint(id))
	if err != nil {
		SendError(400, err.Error(), w)
		return
	}

	SendResponse(user, w)
}