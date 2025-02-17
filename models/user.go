package models

import (
	"boilerplate/services"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	UserName        string    `gorm:"unique" json:"user_name" form:"user_name" validate:"required"`
	Image           string    `json:"image" form:"image"`
	Email           string    `gorm:"unique" json:"email" form:"email" validate:"required,email"`
	Password        string    `json:"-"`
	Salt            string    `json:"-"`
	LastOnline      time.Time `json:"last_online" form:"last_online"`
}

// Not a db object, all the columns that are allowed to be updated by a "normal" way
type UpdateUser struct {
	UserName *string `json:"user_name" form:"user_name"`
	Email    *string `json:"email" form:"email"`
	Online   *bool   `json:"online" form:"online"`
}

type SimpleUser struct {
	ID       uint   `json:"id"`
	UserName string `json:"user_name"`
	Image    string `json:"image"`
}

// Not a db object but as the password shouldn't be returned we want it to be stored apart
type Password struct {
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

// Not a db object, used when updating the password
type UpdatePassword struct {
	OldPassword string `json:"old_password" form:"old_password"`
	Password    string `json:"password" form:"password" validate:"required,min=8"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(user *User, password *Password) []string {
	validator := services.GetValidator()
	errMsgs := make([]string, 0)

	if errs := validator.Validate(*user); len(errs) > 0 && errs[0].Error {
		errMsgs = append(errMsgs, services.FormatErrors(errs)...)
	}
	if errs := validator.Validate(*password); len(errs) > 0 && errs[0].Error {
		errMsgs = append(errMsgs, services.FormatErrors(errs)...)
	}
	if len(errMsgs) > 0 {
		return errMsgs
	}

	user.Salt = services.RandStringBytes(20, true)
	user.Password, _ = HashPassword(password.Password + user.Salt)
	if response := services.DB.Create(user); response.Error != nil {
		errMsgs = append(errMsgs, "error_while_storing" + " user: "+response.Error.Error())
		return errMsgs
	}
	return nil
}

func AreLogInfosCorrect(user *User, password string) (bool, string) {
	if response := services.DB.Where("user_name = ?", user.UserName).Or("email = ?", user.UserName).First(&user); response.Error != nil {
		return false, "invalid_credentials"
	}

	if !CheckPasswordHash(password+user.Salt, user.Password) {
		return false, "invalid_credentials"
	}

	return true, ""
}

func IsPasswordCorrect(password string, id uint) (bool, string) {
	var user User
	if response := services.DB.First(&user, id); response.Error != nil {
		return false, "user_not_found"
	}

	if !CheckPasswordHash(password+user.Salt, user.Password) {
		return false, "wrong_password"
	}

	return true, ""
}

func UpdateUserPassword(password *UpdatePassword, id uint) []string {
	validator := services.GetValidator()
	errMsgs := make([]string, 0)
	var user User

	if response := services.DB.First(&user, id); response.Error != nil {
		errMsgs = append(errMsgs, "user_not_found")
	}
	if errs := validator.Validate(password); len(errs) > 0 && errs[0].Error {
		errMsgs = append(errMsgs, services.FormatErrors(errs)...)
	}
	if len(errMsgs) > 0 {
		return errMsgs
	}

	user.Password, _ = HashPassword(password.Password + user.Salt)
	if response := services.DB.Save(user); response.Error != nil {
		errMsgs = append(errMsgs, "error_while_storing "+" user: "+response.Error.Error())
		return errMsgs
	}
	return nil
}

func GetUser(id uint) (User, error) {
	var user User

	if response := services.DB.
		Select("users.*").
		First(&user, "users.id = ?", id); response.Error != nil {
		return user, response.Error
	}

	if user.Image != "" {
		user.Image, _ = services.GetFile(user.Image)
	}

	return user, nil
}

func GetAllUser(start int, itemsPerPage int) ([]User, error) {
	var users []User

	if response := services.DB.
		Select("users.*").
		Limit(itemsPerPage).
		Offset(start).
		Find(&users); response.Error != nil {
		return users, response.Error
	}

	return users, nil
}

func GetUserByParam(value string, param string) (User, error) {
	var user User

	if response := services.DB.
		Select("users.*").
		First(&user, "users."+param+" = ?", value); response.Error != nil {
		return user, response.Error
	}

	if user.Image != "" {
		user.Image, _ = services.GetFile(user.Image)
	}

	return user, nil
}

func DeleteUser(id uint) error {
	if response := services.DB.Delete(&User{}, id); response.Error != nil {
		return response.Error
	}

	return nil
}