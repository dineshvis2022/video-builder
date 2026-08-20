package main

import (
	"log"
	"video-builder/controllers"
	"video-builder/renderer"
	"video-builder/services"

	"github.com/gin-gonic/gin"
)

func main() {

	// 1. Create FFmpeg renderer
	ffmpegRenderer := renderer.NewFFmpegRenderer(
		"ffmpeg",
		"./generated",
	)

	// 2. Create video service
	videoService := services.NewVideoService(
		ffmpegRenderer,
	)

	// 3. Create handler
	videoHandler := controllers.NewVideoHandler(
		videoService,
	)

	// 4. Create Gin router
	router := gin.Default()

	// 5. Serve generated videos
	router.Static(
		"/generated",
		"./generated",
	)

	// 6. Generate video API
	router.POST(
		"/api/video/generate",
		videoHandler.GenerateVideo,
	)

	// 7. Start server
	log.Println(
		"Video Builder API running on :8080",
	)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
