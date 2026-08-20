package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"video-builder/models"
	"video-builder/services"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	videoService *services.VideoService
}

func NewVideoHandler(videoService *services.VideoService) *VideoHandler {
	return &VideoHandler{
		videoService: videoService,
	}
}

func (h *VideoHandler) GenerateVideo(c *gin.Context) {
	var invitation models.Invitation

	// 1. Decode JSON
	if err := c.ShouldBindJSON(&invitation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid JSON request",
			"error":   err.Error(),
		})
		return
	}

	// 2. Generate video
	outputPath, err := h.videoService.GenerateVideo(invitation)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// 3. Get filename
	filename := filepath.Base(outputPath)

	// 4. Check file exists
	if _, err := os.Stat(outputPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "generated video not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "video generated successfully",
		"videoUrl":  "/generated/" + filename,
		"videoPath": outputPath,
	})
}
