package controllers

import (
	"net/http"

	"neuro-guide-go-service/models"
	"neuro-guide-go-service/services"

	"github.com/gin-gonic/gin"
)

var planService = services.NewPracticePlanService()

// CreatePlanRequest represents a request to create a practice plan
type CreatePlanRequest struct {
	Title string            `json:"title" binding:"required"`
	Days  int               `json:"days" binding:"required"`
	Tasks []models.PlanTask `json:"tasks" binding:"required"`
}

// CreatePlan handles creating a practice plan
func CreatePlan(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan := &models.PracticePlan{
		UserID: userID,
		Title:  req.Title,
		Days:   req.Days,
		Tasks:  req.Tasks,
	}

	if err := planService.CreatePlan(plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plan"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// GetPlan handles getting a practice plan
func GetPlan(c *gin.Context) {
	planID := c.Param("id")
	if planID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plan ID is required"})
		return
	}

	plan, err := planService.GetPlanByID(planID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// GetPlans handles getting all practice plans for a user
func GetPlans(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	plans, err := planService.GetPlansByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get plans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

// DeletePlan handles deleting a practice plan
func DeletePlan(c *gin.Context) {
	planID := c.Param("id")
	if planID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plan ID is required"})
		return
	}

	if err := planService.DeletePlan(planID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan deleted"})
}
