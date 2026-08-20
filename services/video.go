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

	if err := validateAudio(
		invitation.Audio,
	); err != nil {

		return err
	}

	for _, slide := range invitation.Slides {

		// 4. Validate slide duration
		if slide.DisplayDuration <= 0 {
			return errors.New(
				"slide displayDuration must be greater than zero",
			)
		}

		// 4. Validate Transition
		if err := validateTransition(
			slide.Transition,
			slide.DisplayDuration,
		); err != nil {
			return err
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
		if err := validateAnimation(
			element.Animation,
		); err != nil {

			return err
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

func validateTransition(
	transition *models.Transition,
	slideDuration int,
) error {

	if transition == nil {
		return nil
	}

	allowed := map[string]bool{
		"fade":       true,
		"wipeleft":   true,
		"wiperight":  true,
		"wipeup":     true,
		"wipedown":   true,
		"slideleft":  true,
		"slideright": true,
		"slideup":    true,
		"slidedown":  true,
	}

	if !allowed[transition.Type] {

		return fmt.Errorf(
			"unsupported transition type: %s",
			transition.Type,
		)
	}

	if transition.Duration <= 0 {
		return errors.New(
			"transition duration must be greater than zero",
		)
	}

	if transition.Duration >=
		float64(slideDuration) {

		return errors.New(
			"transition duration must be less than slide duration",
		)
	}

	return nil
}

func validateAnimation(
	animation *models.Animation,
) error {

	if animation == nil {
		return nil
	}

	allowed := map[string]bool{
		"fadeIn":     true,
		"fadeOut":    true,
		"slideLeft":  true,
		"slideRight": true,
		"slideUp":    true,
		"slideDown":  true,
		"zoomIn":     true,
		"zoomOut":    true,
	}

	if !allowed[animation.Type] {

		return fmt.Errorf(
			"unsupported animation type: %s",
			animation.Type,
		)
	}

	if animation.Duration < 0 {

		return errors.New(
			"animation duration cannot be negative",
		)
	}

	if animation.Distance < 0 {

		return errors.New(
			"animation distance cannot be negative",
		)
	}

	return nil
}

func validateAudio(
	audio *models.AudioConfig,
) error {

	if audio == nil {
		return nil
	}

	if audio.URL == "" {
		return errors.New(
			"audio url is required",
		)
	}

	if audio.Volume < 0 ||
		audio.Volume > 2 {

		return errors.New(
			"audio volume must be between 0 and 2",
		)
	}

	if audio.BackgroundVolume < 0 ||
		audio.BackgroundVolume > 2 {

		return errors.New(
			"background volume must be between 0 and 2",
		)
	}

	if audio.FadeIn < 0 {
		return errors.New(
			"audio fadeIn cannot be negative",
		)
	}

	if audio.FadeOut < 0 {
		return errors.New(
			"audio fadeOut cannot be negative",
		)
	}

	return nil
}
