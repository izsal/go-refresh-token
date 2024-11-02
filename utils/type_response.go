package utils

import "github.com/izsal/go-refresh-token/models"

type LoginResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type GenericResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type ItemResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Item    models.Item `json:"item"`
}

type ItemsResponse struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Items   []models.Item `json:"items"`
}
