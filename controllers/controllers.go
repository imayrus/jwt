package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/imayrus/jwt-api/database"
	"github.com/imayrus/jwt-api/helpers"
	"github.com/imayrus/jwt-api/models"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		var err models.Error
		err = helpers.SetError(err, "Error in reading payload")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var dbuser models.User
	connection.Where("email = ?", user.Email).First(&dbuser)

	if dbuser.Email != "" {

		var err models.Error
		err = helpers.SetError(err, "Email alreadi in use")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	user.Password, err = helpers.GenerateHashPassword(user.Password)
	if err != nil {
		log.Fatal("Error in hashing password")
	}

	connection.Create(&user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	var authDetails models.Authentication

	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		var err models.Error
		err = helpers.SetError(err, "Error in reading payload")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var authuser models.User
	connection.Where("email=     ?", authDetails.Email).First(&authuser)

	if authuser.Email == "" {
		var err models.Error
		err = helpers.SetError(err, "username is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	checkpassword := helpers.CheckPasswordHash(authDetails.Password, authuser.Password)

	if !checkpassword {
		var err models.Error
		err = helpers.SetError(err, "Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	validToken, err := helpers.GenerateJWT(authuser.Email, authuser.Role, authuser.FirstName)
	if err != nil {
		var err models.Error
		err = helpers.SetError(err, "Failed to generate token")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	var token models.Token
	token.Email = authuser.Email
	token.Role = authuser.Role
	token.FirstName = authuser.FirstName
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)

}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOME PUBLIC INDEX PAGE"))
}

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized"))
		return
	}
	w.Write([]byte("Welcome Admin"))
}

func UserIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "user" {
		w.Write([]byte("Not Authorized"))
		return
	}
	w.Write([]byte("Welcom User"))
}
