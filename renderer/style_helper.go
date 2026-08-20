package renderer

import "video-builder/models"

func getOpacity(element models.Element) float64 {

	if element.Opacity <= 0 {
		return 1
	}

	if element.Opacity > 1 {
		return 1
	}

	return element.Opacity
}

func getRotation(element models.Element) float64 {

	return element.Rotation
}

func getFontSize(element models.Element) int {

	if element.FontSize <= 0 {
		return 60
	}

	return element.FontSize
}

func getFontColor(element models.Element) string {

	if element.FontColor == "" {
		return "white"
	}

	return element.FontColor
}

func getTextAlign(element models.Element) string {

	switch element.Align {

	case "center":
		return "center"

	case "right":
		return "right"

	default:
		return "left"
	}
}
