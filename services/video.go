package services

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"video-builder/models"
	"video-builder/renderer"
)

type VideoService struct {
	renderer *renderer.FFmpegRenderer
}

func NewVideoService(renderer *renderer.FFmpegRenderer) *VideoService {
	return &VideoService{
		renderer: renderer,
	}
}

func (s *VideoService) GenerateVideo(
	invitation models.Invitation,
) (string, error) {

	// 1. Validate request
	if err := validateInvitation(invitation); err != nil {
		return "", err
	}

	// 2. Calculate timeline
	calculateTimeline(&invitation)

	// 3. Render video
	outputPath, err := s.renderer.Render(invitation)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func validateInvitation(invitation models.Invitation) error {

	// 1. Validate project
	if invitation.ProjectID == "" {
		return errors.New("projectId is required")
	}

	// 2. Validate background
	if invitation.Background.URL == "" {
		return errors.New("background video is required")
	}

	// 3. Validate slides
	if len(invitation.Slides) == 0 {
		return errors.New("at least one slide is required")
	}

	for _, slide := range invitation.Slides {

		// 4. Validate slide duration
		if slide.DisplayDuration <= 0 {
			return errors.New(
				"slide displayDuration must be greater than zero",
			)
		}

		// 5. Validate elements
		for _, element := range slide.Elements {

			if err := validateElement(element); err != nil {
				return err
			}
		}
	}

	return nil
}

func calculateTimeline(invitation *models.Invitation) {

	currentTime := 0

	for i := range invitation.Slides {

		invitation.Slides[i].Start = currentTime

		invitation.Slides[i].End =
			currentTime + invitation.Slides[i].DisplayDuration

		currentTime = invitation.Slides[i].End
	}
}

func validateElement(element models.Element) error {

	switch element.Type {

	case "text":

		// 1. Validate text
		if element.Text == "" {
			return errors.New(
				"text element requires text",
			)
		}

		// 2. Validate font size
		if element.FontSize < 0 {
			return errors.New(
				"fontSize cannot be negative",
			)
		}

		// 3. Validate alignment
		if element.Align != "" &&
			element.Align != "left" &&
			element.Align != "center" &&
			element.Align != "right" {

			return errors.New(
				"align must be left, center or right",
			)
		}

		// 4. Validate font
		if element.FontFile != "" {

			if err := validateFontFile(
				element.FontFile,
			); err != nil {
				return err
			}
		}

	case "image":

		// 1. Validate asset path
		if element.AssetPath == "" {
			return errors.New(
				"image element requires assetPath",
			)
		}

		// 2. Validate asset
		if !isSupportedAsset(element.AssetPath) {
			return fmt.Errorf(
				"unsupported image format: %s",
				element.AssetPath,
			)
		}

		// 3. Validate dimensions
		if element.Width <= 0 {
			return errors.New(
				"image width must be greater than zero",
			)
		}

		if element.Height <= 0 {
			return errors.New(
				"image height must be greater than zero",
			)
		}

	case "gif":

		// 1. Validate asset path
		if element.AssetPath == "" {
			return errors.New(
				"gif element requires assetPath",
			)
		}

		// 2. Validate GIF
		if !isGIF(element.AssetPath) {
			return errors.New(
				"gif element requires a .gif file",
			)
		}

		// 3. Validate dimensions
		if element.Width <= 0 {
			return errors.New(
				"gif width must be greater than zero",
			)
		}

		if element.Height <= 0 {
			return errors.New(
				"gif height must be greater than zero",
			)
		}

	default:

		return fmt.Errorf(
			"unsupported element type: %s",
			element.Type,
		)
	}

	// 6. Validate opacity
	if element.Opacity < 0 ||
		element.Opacity > 1 {

		if element.Opacity != 0 {
			return errors.New(
				"opacity must be between 0 and 1",
			)
		}
	}

	return nil
}

func validateFontFile(path string) error {

	// 1. Check extension
	ext := strings.ToLower(
		filepath.Ext(path),
	)

	switch ext {
	case ".ttf", ".otf":
	default:
		return fmt.Errorf(
			"unsupported font format: %s",
			ext,
		)
	}

	return nil
}

func isSupportedAsset(path string) bool {

	ext := strings.ToLower(
		filepath.Ext(path),
	)

	switch ext {
	case ".jpg",
		".jpeg",
		".png",
		".webp":

		return true
	default:
		return false
	}
}

func isGIF(path string) bool {

	return strings.EqualFold(
		filepath.Ext(path),
		".gif",
	)
}
