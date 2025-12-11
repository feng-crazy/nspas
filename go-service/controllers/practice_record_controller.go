package controllers

import (
	"net/http"
	"time"

	"neuro-guide-go-service/models"
	"neuro-guide-go-service/services"

	"github.com/gin-gonic/gin"
)

var recordService = services.NewPracticeRecordService()

// CreateRecordRequest represents a request to create a practice record
type CreateRecordRequest struct {
	PlanID         string    `json:"plan_id" binding:"required"`
	Date           time.Time `json:"date" binding:"required"`
	CompletedTasks []string  `json:"completed_tasks"`
	Reflection     string    `json:"reflection"`
}

// CreateRecord handles creating a practice record
func CreateRecord(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record := &models.PracticeRecord{
		UserID:         userID,
		PlanID:         req.PlanID,
		Date:           req.Date,
		CompletedTasks: req.CompletedTasks,
		Reflection:     req.Reflection,
	}

	if err := recordService.CreateRecord(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// UpdateRecordRequest represents a request to update a practice record
type UpdateRecordRequest struct {
	CompletedTasks []string `json:"completed_tasks"`
	Reflection     string   `json:"reflection"`
}

// UpdateRecord handles updating a practice record
func UpdateRecord(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID is required"})
		return
	}

	var req UpdateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record, err := recordService.GetRecordByID(recordID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	if req.CompletedTasks != nil {
		record.CompletedTasks = req.CompletedTasks
	}
	if req.Reflection != "" {
		record.Reflection = req.Reflection
	}

	if err := recordService.UpdateRecord(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// GetRecord handles getting a practice record
func GetRecord(c *gin.Context) {
	recordID := c.Param("id")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record ID is required"})
		return
	}

	record, err := recordService.GetRecordByID(recordID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// GetRecords handles getting all practice records
func GetRecords(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	planID := c.Query("plan_id")
	if planID != "" {
		records, err := recordService.GetRecordsByPlanID(planID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get records"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"records": records})
		return
	}

	records, err := recordService.GetRecordsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"records": records})
}
