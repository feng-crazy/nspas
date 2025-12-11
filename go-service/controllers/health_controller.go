package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheckResponse defines the structure for health check responses
type HealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HealthCheck handles health check requests
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthCheckResponse{
		Status:  "ok",
		Message: "Go service is running",
	})
}
