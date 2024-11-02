package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/izsal/go-refresh-token/database"
	"github.com/izsal/go-refresh-token/logger"
	"github.com/izsal/go-refresh-token/models"
	"github.com/izsal/go-refresh-token/utils"

	"gorm.io/gorm"
)

var logFile *os.File

func init() {
	logFile = logger.InitLogger()
	log.SetOutput(logFile) // Set logger to write to logFile
}

func logUserAction(username, action, status string) {
	log.Printf("User: %s | Action: %s | Status: %s\n", username, action, status)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		logUserAction(user.Username, "Register", "Failed - Invalid input")
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		logUserAction(user.Username, "Register", "Failed - Missing credentials")
		return
	}

	var existingUser models.User
	result := database.DB.Where("username = ?", user.Username).First(&existingUser)
	if result.Error == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		logUserAction(user.Username, "Register", "Failed - Username already exists")
		return
	} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Error checking username", http.StatusInternalServerError)
		logUserAction(user.Username, "Register", "Failed - Database error")
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		logUserAction(user.Username, "Register", "Failed - Password hashing error")
		return
	}
	user.Password = hashedPassword

	result = database.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		logUserAction(user.Username, "Register", "Failed - User creation error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"success": "true", "message": "User created successfully"})
	logUserAction(user.Username, "Register", "Success")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		logUserAction(user.Username, "Login", "Failed - Invalid input")
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		logUserAction(user.Username, "Login", "Failed - Missing credentials")
		return
	}

	var dbUser models.User
	result := database.DB.Where("username = ?", user.Username).First(&dbUser)
	if result.Error == gorm.ErrRecordNotFound || !utils.CheckPasswordHash(user.Password, dbUser.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		logUserAction(user.Username, "Login", "Failed - Invalid credentials")
		return
	}

	token, err := utils.GenerateJWT(dbUser.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		logUserAction(user.Username, "Login", "Failed - Token generation error")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(dbUser.Username)
	if err != nil {
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		logUserAction(user.Username, "Login", "Failed - Refresh token generation error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := utils.LoginResponse{
		Status:       "success",
		Message:      "Login successful",
		Token:        token,
		RefreshToken: refreshToken,
	}
	json.NewEncoder(w).Encode(response)
	logUserAction(user.Username, "Login", "Success")
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		logUserAction("", "RefreshToken", "Failed - Invalid input")
		return
	}

	refreshToken := requestBody["refreshToken"]
	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		logUserAction("", "RefreshToken", "Failed - Invalid refresh token")
		return
	}

	newToken, err := utils.GenerateJWT(claims.Username)
	if err != nil {
		http.Error(w, "Error generating new token", http.StatusInternalServerError)
		logUserAction(claims.Username, "RefreshToken", "Failed - Token generation error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": newToken})
	logUserAction(claims.Username, "RefreshToken", "Success")
}
