package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/izsal/go-refresh-token/database"
	"github.com/izsal/go-refresh-token/models"
	"github.com/izsal/go-refresh-token/utils"

	"gorm.io/gorm"
)


func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate input
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Check if the username already exists
	var existingUser models.User
	result := database.DB.Where("username = ?", user.Username).First(&existingUser)
	if result.Error == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Error checking username", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	result = database.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Create response message
	response := map[string]string{"success": "true", "message": "User created successfully"}
	json.NewEncoder(w).Encode(response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate input
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	var dbUser models.User
	result := database.DB.Where("username = ?", user.Username).First(&dbUser)
	if result.Error == gorm.ErrRecordNotFound || !utils.CheckPasswordHash(user.Password, dbUser.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(dbUser.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		log.Println("Error generating JWT:", err)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(dbUser.Username)
	if err != nil {
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		log.Println("Error generating refresh token:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := utils.LoginResponse{
		Status:       "success",
		Message:      "Login successful",
		Token:        token,
		RefreshToken: refreshToken,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	refreshToken := requestBody["refreshToken"]
	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	newToken, err := utils.GenerateJWT(claims.Username)
	if err != nil {
		http.Error(w, "Error generating new token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{
		"token": newToken,
	})
	if err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
