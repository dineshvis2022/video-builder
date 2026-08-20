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

	currentVideo := "[0:v]"

	filterCount := 0

	for _, slide := range invitation.Slides {

		for _, element := range slide.Elements {

			switch element.Type {

			case "text":

				outputLabel := fmt.Sprintf(
					"[v%d]",
					filterCount,
				)

				filter := buildTextFilter(
					currentVideo,
					outputLabel,
					element,
					slide.Start,
					slide.End,
				)

				b.filters = append(
					b.filters,
					filter,
				)

				currentVideo = outputLabel

				filterCount++

			case "image":

				inputLabel := fmt.Sprintf(
					"[%d:v]",
					b.inputIndex,
				)

				outputLabel := fmt.Sprintf(
					"[v%d]",
					filterCount,
				)

				filter := buildImageFilter(
					currentVideo,
					inputLabel,
					outputLabel,
					element,
					slide.Start,
					slide.End,
				)

				b.filters = append(
					b.filters,
					filter,
				)

				currentVideo = outputLabel

				b.inputIndex++
				filterCount++

			case "gif":

				inputLabel := fmt.Sprintf(
					"[%d:v]",
					b.inputIndex,
				)

				outputLabel := fmt.Sprintf(
					"[v%d]",
					filterCount,
				)

				filter := buildGIFFilter(
					currentVideo,
					inputLabel,
					outputLabel,
					element,
					slide.Start,
					slide.End,
				)

				b.filters = append(
					b.filters,
					filter,
				)

				currentVideo = outputLabel

				b.inputIndex++
				filterCount++

			default:

				return "", 0, fmt.Errorf(
					"unsupported element type: %s",
					element.Type,
				)
			}
		}
	}

	b.filters = append(
		b.filters,
		fmt.Sprintf(
			"%snull[vout]",
			currentVideo,
		),
	)

	return strings.Join(b.filters, ";"), b.inputIndex - 1, nil
}

func buildTextFilter(
	inputLabel string,
	outputLabel string,
	element models.Element,
	start int,
	end int,
) string {

	fontSize := getFontSize(element)
	fontColor := getFontColor(element)
	opacity := getOpacity(element)
	align := getTextAlign(element)

	text := escapeFFmpegText(element.Text)

	fontColorWithOpacity := fmt.Sprintf(
		"%s@%f",
		fontColor,
		opacity,
	)

	xExpression := buildTextXExpression(
		element.X,
		align,
	)

	fontFile := ""

	if element.FontFile != "" {

		escapedFontFile :=
			escapeFFmpegPath(element.FontFile)

		fontFile = fmt.Sprintf(
			"fontfile='%s':",
			escapedFontFile,
		)
	}

	return fmt.Sprintf(
		"%sdrawtext="+
			"%stext='%s':"+
			"x=%s:"+
			"y=%d:"+
			"fontsize=%d:"+
			"fontcolor=%s:"+
			"enable='between(t,%d,%d)'%s",

		inputLabel,

		fontFile,

		text,

		xExpression,

		element.Y,

		fontSize,

		fontColorWithOpacity,

		start,

		end,

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

func buildGIFFilter(
	videoLabel string,
	gifLabel string,
	outputLabel string,
	element models.Element,
	start int,
	end int,
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

	return fmt.Sprintf(
		"%s%s[%s];"+
			"%s[%s]overlay=%d:%d:"+
			"enable='between(t,%d,%d)'%s",

		gifLabel,

		gifFilters,

		gifLabelName,

		videoLabel,

		gifLabelName,

		element.X,

		element.Y,

		start,

		end,

		outputLabel,
	)
}

func buildTextXExpression(
	x int,
	align string,
) string {

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
