package controllers

import (
	"encoding/json"
	"fmt"
	"go-plate/models"
	"go-plate/services"
	"go-plate/translations"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(400, translations.T(r.Context().Value(services.LangKey).(string), "failed_to_read_body"), w)
		return
	}
	
	user := new(models.User)
	user.LastOnline = time.Now()

	if err := json.Unmarshal(body, &user); err != nil {
		SendError(422, translations.T(r.Context().Value(services.LangKey).(string), "fail_to_parse_body"), w)
		return
	}

	password := new(models.Password)
	if err := json.Unmarshal(body, &password); err != nil {
		fmt.Println(err)
		SendError(422, translations.T(r.Context().Value(services.LangKey).(string), "fail_to_retrieve_password"), w)
		return
	}

	errs := models.CreateUser(user, password)
	if errs != nil {
		SendError(400, errs[0], w)
		return
	}

	if _, err := services.SendEmail(user.Email, translations.T(r.Context().Value(services.LangKey).(string), "welcome"), translations.T(r.Context().Value(services.LangKey).(string), "welcome_message")); err != nil {
		SendError(400, translations.T(r.Context().Value(services.LangKey).(string), "fail_to_send_email"), w)
		return
	}

	SendResponse(user, w)
}

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(400, translations.T(r.Context().Value(services.LangKey).(string), "failed_to_read_body"), w)
		return
	}
	
	user := new(models.User)
	if err := json.Unmarshal(body, &user); err != nil {
		SendError(422, translations.T(r.Context().Value(services.LangKey).(string), "fail_to_parse_body"), w)
		return
	}

	password := new(models.Password)
	if err := json.Unmarshal(body, &password); err != nil {
		fmt.Println(err)
		SendError(422, translations.T(r.Context().Value(services.LangKey).(string), "fail_to_retrieve_password"), w)
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
		SendError(400, translations.T(r.Context().Value(services.LangKey).(string), "fail_to_retrieve_users"), w)
		return
	}

	SendResponse(users, w)
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	user, err := models.GetUser(r.Context().Value(services.UserIDKey).(uint))
	if err != nil {
		SendError(400, err.Error(), w)
	}

	SendResponse(user, w)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, translations.T(r.Context().Value(services.LangKey).(string), "invalid_id_format"), http.StatusBadRequest)
		return
	}

	user, err := models.GetUser(uint(id))
	if err != nil {
		SendError(400, err.Error(), w)
		return
	}

	SendResponse(user, w)
}

func UploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	ext := filepath.Ext(handler.Filename)
	if !allowedExtensions[ext] {
		http.Error(w, "Invalid file format", http.StatusBadRequest)
		return
	}

	filename, err, _ := services.PutImage(&file, ext, "profile", r.Context().Value(services.UserIDKey).(uint))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error uploading the file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "Upload successful", "image_url": "%s"}`, filename)))
}