package renderer

import (
	"fmt"
	"strings"
	"video-builder/models"
)

func getAnimationDuration(
	animation *models.Animation,
) float64 {

	if animation == nil {
		return 0
	}

	if animation.Duration <= 0 {
		return 1
	}

	return animation.Duration
}

func getAnimationDistance(
	animation *models.Animation,
) int {

	if animation == nil {
		return 500
	}

	if animation.Distance <= 0 {
		return 500
	}

	return animation.Distance
}

func buildAnimationProgress(
	animation *models.Animation,
) string {

	duration := getAnimationDuration(animation)

	return fmt.Sprintf(
		"min(t/%f,1)",
		duration,
	)
}

func buildFadeInExpression(
	animation *models.Animation,
	opacity float64,
) string {

	progress := buildAnimationProgress(
		animation,
	)

	return fmt.Sprintf(
		"%f*%s",
		opacity,
		progress,
	)
}

func buildFadeOutExpression(
	animation *models.Animation,
	opacity float64,
	slideDuration int,
) string {

	duration := getAnimationDuration(
		animation,
	)

	start := float64(slideDuration) - duration

	return fmt.Sprintf(
		"%f*(1-max((t-%f)/%f,0))",
		opacity,
		start,
		duration,
	)
}

func buildSlideXExpression(
	animation *models.Animation,
	baseX string,
) string {

	progress := buildAnimationProgress(
		animation,
	)

	distance := getAnimationDistance(
		animation,
	)

	switch animation.Type {

	case "slideLeft":

		return fmt.Sprintf(
			"%s-%d*(1-%s)",
			baseX,
			distance,
			progress,
		)

	case "slideRight":

		return fmt.Sprintf(
			"%s+%d*(1-%s)",
			baseX,
			distance,
			progress,
		)
	}

	return baseX
}

func buildSlideYExpression(
	animation *models.Animation,
	baseY string,
) string {

	progress := buildAnimationProgress(
		animation,
	)

	distance := getAnimationDistance(
		animation,
	)

	switch animation.Type {

	case "slideUp":

		return fmt.Sprintf(
			"%s-%d*(1-%s)",
			baseY,
			distance,
			progress,
		)

	case "slideDown":

		return fmt.Sprintf(
			"%s+%d*(1-%s)",
			baseY,
			distance,
			progress,
		)
	}

	return baseY
}

func buildGIFFilter(
	videoLabel string,
	gifLabel string,
	outputLabel string,
	element models.Element,
	slideDuration int,
) string {

	gifLabelName :=
		strings.Trim(
			outputLabel,
			"[]",
		) + "_gif"

	opacity := getOpacity(element)

	rotation := getRotation(element)

	scaleFilter := fmt.Sprintf(
		"scale=%d:%d:force_original_aspect_ratio=decrease",
		element.Width,
		element.Height,
	)

	gifFilters :=
		"format=rgba," + scaleFilter

	if rotation != 0 {

		gifFilters += fmt.Sprintf(
			",rotate=%f*PI/180:"+
				"ow=rotw(iw):"+
				"oh=roth(ih):"+
				"c=none",
			rotation,
		)
	}

	gifFilters += fmt.Sprintf(
		",colorchannelmixer=aa=%f",
		opacity,
	)
	if element.Animation != nil {

		switch element.Animation.Type {

		case "fadeIn":

			duration :=
				getAnimationDuration(
					element.Animation,
				)

			gifFilters += fmt.Sprintf(
				",fade=t=in:st=0:d=%f:alpha=1",
				duration,
			)

		case "fadeOut":

			duration :=
				getAnimationDuration(
					element.Animation,
				)

			start :=
				float64(slideDuration) - duration

			gifFilters += fmt.Sprintf(
				",fade=t=out:st=%f:d=%f:alpha=1",
				start,
				duration,
			)
		}
	}

	baseX := fmt.Sprintf(
		"%d",
		element.X,
	)

	baseY := fmt.Sprintf(
		"%d",
		element.Y,
	)

	xExpression := baseX
	yExpression := baseY

	if element.Animation != nil {

		switch element.Animation.Type {

		case "slideLeft",
			"slideRight":

			xExpression =
				buildSlideXExpression(
					element.Animation,
					baseX,
				)

		case "slideUp",
			"slideDown":

			yExpression =
				buildSlideYExpression(
					element.Animation,
					baseY,
				)
		}
	}
	return fmt.Sprintf(
		"%s%s[%s];"+
			"%s[%s]overlay=x='%s':y='%s'%s",

		gifLabel,

		gifFilters,

		gifLabelName,

		videoLabel,

		gifLabelName,

		xExpression,

		yExpression,

		outputLabel,
	)
}
