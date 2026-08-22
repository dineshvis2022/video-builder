package renderer

import (
	"fmt"
	"os"
	"strings"
	"video-builder/models"
)

type FilterBuilder struct {
	inputIndex int
	filters    []string
}

func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{
		inputIndex: 1,
		filters:    make([]string, 0),
	}
}

func (b *FilterBuilder) Build(
	invitation models.Invitation,
) (string, int, error) {

	if len(invitation.Slides) == 0 {
		return "", 0, fmt.Errorf(
			"no slides found",
		)
	}

	// 1. Split background into slide streams
	splitLabels := make([]string, 0)
	// splitExpression := fmt.Sprintf(
	// 	"[0:v]split=%d%s",
	// 	len(invitation.Slides),
	// 	strings.Join(splitLabels, ""),
	// )

	for i := range invitation.Slides {

		splitLabels = append(
			splitLabels,
			fmt.Sprintf(
				"[bg%d]",
				i,
			),
		)
	}

	b.filters = append(
		b.filters,
		fmt.Sprintf(
			"[0:v]split=%d%s",
			len(invitation.Slides),
			strings.Join(
				splitLabels,
				"",
			),
		),
	)

	slideLabels := make(
		[]string,
		0,
		len(invitation.Slides),
	)

	// 2. Build every slide
	for i, slide := range invitation.Slides {

		inputLabel := fmt.Sprintf(
			"[bg%d]",
			i,
		)

		outputLabel := fmt.Sprintf(
			"[slide%d]",
			i,
		)

		trimFilter := fmt.Sprintf(
			"%strim=start=%d:end=%d,setpts=PTS-STARTPTS%s",
			inputLabel,
			slide.Start,
			slide.End,
			outputLabel,
		)

		b.filters = append(
			b.filters,
			trimFilter,
		)

		currentVideo := outputLabel

		// 3. Add slide elements
		for _, element := range slide.Elements {

			switch element.Type {

			case "text":

				nextLabel := fmt.Sprintf(
					"[s%d_%d]",
					i,
					len(b.filters),
				)

				filter := buildTextFilter(
					currentVideo,
					nextLabel,
					element,
					slide.DisplayDuration,
				)

				b.filters = append(
					b.filters,
					filter,
				)

				currentVideo = nextLabel

			case "image":

				inputLabel := fmt.Sprintf(
					"[%d:v]",
					b.inputIndex,
				)

				nextLabel := fmt.Sprintf(
					"[s%d_%d]",
					i,
					len(b.filters),
				)

				filter := buildImageFilter(
					currentVideo,
					inputLabel,
					nextLabel,
					element,
					slide.DisplayDuration,
				)

				b.filters = append(
					b.filters,
					filter,
				)

				currentVideo = nextLabel

				b.inputIndex++

			case "gif":

				inputLabel := fmt.Sprintf(
					"[%d:v]",
					b.inputIndex,
				)

				nextLabel := fmt.Sprintf(
					"[s%d_%d]",
					i,
					len(b.filters),
				)

				filter := buildGIFFilter(
					currentVideo,
					inputLabel,
					nextLabel,
					element,
					slide.DisplayDuration,
				)

				b.filters = append(
					b.filters,
					filter,
				)

				currentVideo = nextLabel

				b.inputIndex++

			default:

				return "", 0, fmt.Errorf(
					"unsupported element type: %s",
					element.Type,
				)
			}
		}

		slideLabels = append(
			slideLabels,
			currentVideo,
		)
	}

	// 4. Build slide transitions
	currentVideo := slideLabels[0]

	currentDuration :=
		float64(invitation.Slides[0].DisplayDuration)

	for i := 1; i < len(slideLabels); i++ {

		transition := invitation.Slides[i].Transition

		if transition == nil {

			transition = &models.Transition{
				Type:     "fade",
				Duration: 0,
			}
		}

		duration := transition.Duration

		if duration <= 0 {

			nextVideo := fmt.Sprintf(
				"[concat%d]",
				i,
			)

			b.filters = append(
				b.filters,
				fmt.Sprintf(
					"%s%sconcat=n=2:v=1:a=0%s",
					currentVideo,
					slideLabels[i],
					nextVideo,
				),
			)

			currentVideo = nextVideo

			currentDuration +=
				float64(
					invitation.Slides[i].DisplayDuration,
				)

			continue
		}

		offset :=
			currentDuration - duration

		nextVideo := fmt.Sprintf(
			"[transition%d]",
			i,
		)

		b.filters = append(
			b.filters,
			fmt.Sprintf(
				"%s%sxfade=transition=%s:duration=%f:offset=%f%s",
				currentVideo,
				slideLabels[i],
				transition.Type,
				duration,
				offset,
				nextVideo,
			),
		)

		currentVideo = nextVideo

		currentDuration =
			currentDuration +
				float64(
					invitation.Slides[i].DisplayDuration,
				) -
				duration
	}

	// 5. Final output
	b.filters = append(
		b.filters,
		fmt.Sprintf(
			"%sformat=yuv420p[vout]",
			currentVideo,
		),
	)

	return strings.Join(
		b.filters,
		";",
	), b.inputIndex - 1, nil
}

func buildTextFilter(
	inputLabel string,
	outputLabel string,
	element models.Element,
	slideDuration int,
) string {

	fontSize := getFontSize(element)
	fontColor := getFontColor(element)
	opacity := getOpacity(element)
	align := getTextAlign(element)

	text := escapeFFmpegText(
		element.Text,
	)

	fontColorWithOpacity := fontColor

	baseX := buildTextXExpression(
		element.X,
		align,
	)

	baseY := fmt.Sprintf(
		"%d",
		element.Y,
	)

	xExpression := baseX
	yExpression := baseY

	alphaExpression := fmt.Sprintf(
		"%f",
		opacity,
	)

	if element.Animation != nil {

		switch element.Animation.Type {

		case "fadeIn":

			alphaExpression =
				buildFadeInExpression(
					element.Animation,
					opacity,
				)

		case "fadeOut":

			alphaExpression =
				buildFadeOutExpression(
					element.Animation,
					opacity,
					slideDuration,
				)

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

	fontFile := ""

	if element.FontFile != "" {

		escapedFontFile :=
			escapeFFmpegPath(
				element.FontFile,
			)

		fontFile = fmt.Sprintf(
			"fontfile='%s':",
			escapedFontFile,
		)
	}

	return fmt.Sprintf(
		"%sdrawtext="+
			"%stext='%s':"+
			"x='%s':"+
			"y='%s':"+
			"fontsize=%d:"+
			"fontcolor=%s:"+
			"alpha='%s'%s",

		inputLabel,

		fontFile,

		text,

		xExpression,

		yExpression,

		fontSize,

		fontColorWithOpacity,

		alphaExpression,

		outputLabel,
	)
}

func validateFontFile(
	element models.Element,
) error {

	if element.FontFile == "" {
		return nil
	}

	info, err := os.Stat(
		element.FontFile,
	)

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(
				"font file not found: %s",
				element.FontFile,
			)
		}

		return err
	}

	if info.IsDir() {
		return fmt.Errorf(
			"font path is a directory: %s",
			element.FontFile,
		)
	}

	return nil
}

func escapeFFmpegPath(path string) string {

	path = strings.ReplaceAll(
		path,
		"\\",
		"\\\\",
	)

	path = strings.ReplaceAll(
		path,
		"'",
		"\\'",
	)

	path = strings.ReplaceAll(
		path,
		":",
		"\\:",
	)

	return path
}

func buildTextXExpression(x int, align string) string {

	switch align {

	case "center":

		return fmt.Sprintf(
			"%d-text_w/2",
			x,
		)

	case "right":

		return fmt.Sprintf(
			"%d-text_w",
			x,
		)

	default:

		return fmt.Sprintf(
			"%d",
			x,
		)
	}
}

func buildImageFilter(
	videoLabel string,
	imageLabel string,
	outputLabel string,
	element models.Element,
	slideDuration int,
) string {

	imageLabelName :=
		strings.Trim(
			outputLabel,
			"[]",
		) + "_img"

	opacity := getOpacity(element)

	rotation := getRotation(element)

	scaleFilter := fmt.Sprintf(
		"scale=%d:%d:force_original_aspect_ratio=decrease",
		element.Width,
		element.Height,
	)

	imageFilters := scaleFilter

	if rotation != 0 {

		imageFilters += fmt.Sprintf(
			",rotate=%f*PI/180:"+
				"ow=rotw(iw):"+
				"oh=roth(ih):"+
				"c=none",
			rotation,
		)
	}

	imageFilters +=
		",format=rgba"

	imageFilters += fmt.Sprintf(
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

			imageFilters += fmt.Sprintf(
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

			imageFilters += fmt.Sprintf(
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
			"%s[%s]overlay=x='%s':y='%s':shortest=1%s",

		imageLabel,
		imageFilters,
		imageLabelName,

		videoLabel,
		imageLabelName,

		xExpression,
		yExpression,

		outputLabel,
	)
}
