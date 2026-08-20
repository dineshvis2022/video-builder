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

func (s *VideoService) GenerateVideo(invitation models.Invitation) (string, error) {
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
	if invitation.ProjectID == "" {
		return errors.New("projectId is required")
	}

	if invitation.Background.URL == "" {
		return errors.New("background video is required")
	}

	if len(invitation.Slides) == 0 {
		return errors.New("at least one slide is required")
	}

	for _, slide := range invitation.Slides {

		if slide.DisplayDuration <= 0 {
			return errors.New(
				"slide displayDuration must be greater than zero",
			)
		}

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

func isSupportedAsset(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		return true
	default:
		return false
	}
}

func isGIF(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".gif")
}

func validateElement(element models.Element) error {

	switch element.Type {

	case "text":

		if element.Text == "" {
			return errors.New(
				"text element requires text",
			)
		}

		if element.FontSize < 0 {
			return errors.New(
				"fontSize cannot be negative",
			)
		}

	case "image":

		if element.AssetPath == "" {
			return errors.New(
				"image element requires assetPath",
			)
		}

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

		if element.AssetPath == "" {
			return errors.New(
				"gif element requires assetPath",
			)
		}

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
